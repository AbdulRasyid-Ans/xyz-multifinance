package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
)

type TransactionUsecase interface {
	GetRemainingPayment(ctx context.Context, req TransactionRequest) (response RemainingPaymentResponse, err error)
	CreateTransaction(ctx context.Context, req TransactionRequest) (response GetTransactionResponse, err error)
}

type transactionUsecase struct {
	transactionRepo   repository.TransactionRepository
	loanRepo          repository.LoanRepository
	consumerLimitRepo repository.ConsumerLimitRepository
	consumerRepo      repository.ConsumerRepository
	ctxTimeout        time.Duration
	sync.Mutex
}

type (
	TransactionRequest struct {
		ConsumerID      int64  `json:"consumer_id" query:"consumer_id"`
		LoanID          int64  `json:"loan_id" query:"loan_id"`
		TramsactionType string `json:"transaction_type" query:"transaction_type"`
	}

	RemainingPaymentResponse struct {
		ContractNumber          string    `json:"contract_number"`
		Installment             int32     `json:"installment,omitempty"`
		Tenure                  int16     `json:"tenure"`
		DueDate                 time.Time `json:"due_date"`
		PaidLoanAmount          float64   `json:"paid_loan_amount"`
		PaidInterestAmount      float64   `json:"paid_interest_amount"`
		RemainingLoanAmount     float64   `json:"remaining_loan_amount"`
		RemainingInterestAmount float64   `json:"remaining_interest_amount"`
		TotalRemainingAmount    float64   `json:"total_remaining_amount"`
	}

	GetTransactionResponse struct {
		ID          int64   `json:"id"`
		ConsumerID  int64   `json:"consumer_id"`
		LoanID      int64   `json:"loan_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
)

var (
	validTransactionType = map[string]bool{
		"installment": true,
		"full":        true,
	}
)

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	loanRepo repository.LoanRepository,
	consumerLimitRepo repository.ConsumerLimitRepository,
	consumerRepo repository.ConsumerRepository,
	timeout time.Duration,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo:   transactionRepo,
		loanRepo:          loanRepo,
		consumerLimitRepo: consumerLimitRepo,
		consumerRepo:      consumerRepo,
		ctxTimeout:        timeout,
	}
}

func (uc *transactionUsecase) GetRemainingPayment(ctx context.Context, req TransactionRequest) (response RemainingPaymentResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	_, ok := validTransactionType[req.TramsactionType]
	if !ok {
		return response, errors.New("transaction_type must be installment or full")
	}

	consumer, err := uc.consumerRepo.GetConsumerByID(ctx, req.ConsumerID)
	if err != nil {
		return response, err
	}
	if consumer.ID == 0 {
		return response, errors.New("consumer not found")
	}

	loan, err := uc.loanRepo.GetLoanByID(ctx, req.LoanID)
	if err != nil {
		return response, err
	}

	if loan.ID == 0 {
		return response, errors.New("loan not found")
	}

	if loan.ConsumerID != consumer.ID {
		return response, errors.New("loan not belong to consumer")
	}

	if loan.LoanStatus == "finish" {
		return response, errors.New("loan already finished")
	}

	consumerLimit, err := uc.consumerLimitRepo.GetConsumerLimitByID(ctx, loan.ConsumerLimitID)
	if err != nil {
		return response, err
	}

	if consumerLimit.ID == 0 {
		return response, errors.New("consumer limit not found")
	}

	if req.TramsactionType == "full" {
		response.ContractNumber = loan.ContractNumber
		response.Tenure = consumerLimit.Tenure
		response.DueDate = loan.DueDate
		response.PaidLoanAmount = loan.PaidLoanAmount
		response.PaidInterestAmount = loan.PaidInterestAmount

		response.RemainingLoanAmount = loan.LoanAmount - loan.PaidLoanAmount
		response.RemainingInterestAmount = loan.InterestAmount - loan.PaidInterestAmount
		response.TotalRemainingAmount = response.RemainingLoanAmount + response.RemainingInterestAmount
	} else {
		loanAmount := loan.LoanAmount / float64(consumerLimit.Tenure)
		interestAmount := loan.InterestAmount / float64(consumerLimit.Tenure)

		response.ContractNumber = loan.ContractNumber
		response.Installment = loan.Installment + 1
		response.Tenure = consumerLimit.Tenure
		response.DueDate = loan.DueDate
		response.PaidLoanAmount = loan.PaidLoanAmount
		response.PaidInterestAmount = loan.PaidInterestAmount

		response.RemainingLoanAmount = loanAmount
		response.RemainingInterestAmount = interestAmount
		response.TotalRemainingAmount = response.RemainingLoanAmount + response.RemainingInterestAmount
	}

	return response, nil
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, req TransactionRequest) (response GetTransactionResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	uc.Lock()
	defer uc.Unlock()

	remainingPayemnt, err := uc.GetRemainingPayment(ctx, req)
	if err != nil {
		return response, err
	}

	tx, err := uc.transactionRepo.BeginTx(ctx)
	if err != nil {
		return response, err
	}

	now := time.Now()
	loanStatus := "on_going"
	installment := remainingPayemnt.Installment

	if now.After(remainingPayemnt.DueDate) {
		loanStatus = "late"
	}

	if req.TramsactionType == "full" || remainingPayemnt.Installment == int32(remainingPayemnt.Tenure) {
		loanStatus = "finish"
		installment = int32(remainingPayemnt.Tenure)
	}

	transaction := repository.Transaction{
		ConsumerID:  req.ConsumerID,
		LoanID:      req.LoanID,
		Amount:      remainingPayemnt.TotalRemainingAmount,
		Description: fmt.Sprintf("Payment for loan %s", remainingPayemnt.ContractNumber),
	}

	transactionID, err := uc.transactionRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		uc.transactionRepo.RollbackTx(ctx, tx)
		return response, err
	}

	loan := repository.UpdateLoanRequest{
		ID:                 req.LoanID,
		PaidLoanAmount:     remainingPayemnt.PaidLoanAmount + remainingPayemnt.RemainingLoanAmount,
		PaidInterestAmount: remainingPayemnt.PaidInterestAmount + remainingPayemnt.RemainingInterestAmount,
		LoanStatus:         loanStatus,
		Installment:        installment,
	}

	err = uc.loanRepo.UpdateLoan(ctx, loan, tx)
	if err != nil {
		uc.transactionRepo.RollbackTx(ctx, tx)
		return response, err
	}

	err = uc.transactionRepo.CommitTx(ctx, tx)
	if err != nil {
		return response, err
	}

	response.ID = transactionID
	response.ConsumerID = req.ConsumerID
	response.LoanID = req.LoanID
	response.Amount = transaction.Amount
	response.Description = transaction.Description

	return response, nil
}
