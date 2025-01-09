package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
)

type LoanUsecase interface {
	CreateLoan(ctx context.Context, req CreateLoanRequest) (response LoanResponse, err error)
	GetLoanByID(ctx context.Context, loanID int64) (response LoanResponse, err error)
	GetLoanByConsumerID(ctx context.Context, consumerID int64) (response []LoanResponse, err error)
	DeleteLoanByID(ctx context.Context, loanID int64) (err error)
}

type loanUsecase struct {
	loanRepo          repository.LoanRepository
	consumerLimitRepo repository.ConsumerLimitRepository
	consumerRepo      repository.ConsumerRepository
	merchantRepo      repository.MerchantRepository
	ctxTimeout        time.Duration
}

type (
	CreateLoanRequest struct {
		ConsumerID   int64   `json:"consumer_id"`
		MerchantID   int64   `json:"merchant_id"`
		Tenure       int16   `json:"tenure"`
		LoanAmount   float64 `json:"loan_amount"`
		InterestRate float64 `json:"interest_rate"`
		AssetName    string  `json:"asset_name"`
	}

	LoanResponse struct {
		ID              int64   `json:"id"`
		ConsumerID      int64   `json:"consumer_id"`
		MerchantID      int64   `json:"merchant_id"`
		ConsumerLimitID int64   `json:"consumer_limit_id"`
		LoanAmount      float64 `json:"loan_amount"`
		ContractNumber  string  `json:"contract_number"`
		InterestRate    float64 `json:"interest_rate"`
		InterestAmount  float64 `json:"interest_amount"`
		LoanStatus      string  `json:"loan_status"`
		DueDate         string  `json:"due_date"`
		Installment     int32   `json:"installment"`
		AssetName       string  `json:"asset_name"`
		CreatedAt       string  `json:"created_at,omitempty"`
		UpdatedAt       string  `json:"updated_at,omitempty"`
	}
)

func NewLoanUsecase(
	loanRepo repository.LoanRepository,
	consumerLimitRepo repository.ConsumerLimitRepository,
	consumerRepo repository.ConsumerRepository,
	merchantRepo repository.MerchantRepository,
	timeout time.Duration,
) LoanUsecase {
	return &loanUsecase{
		loanRepo:          loanRepo,
		consumerLimitRepo: consumerLimitRepo,
		consumerRepo:      consumerRepo,
		merchantRepo:      merchantRepo,
		ctxTimeout:        timeout,
	}
}

func (uc *loanUsecase) CreateLoan(ctx context.Context, req CreateLoanRequest) (response LoanResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	consumer, err := uc.consumerRepo.GetConsumerByID(ctx, req.ConsumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanUsecase][CreateLoan] while get consumer by ID, Err: %+v", err))
		return response, err
	}
	if consumer.ID == 0 {
		return response, errors.New("consumer not found")
	}

	merchant, err := uc.merchantRepo.GetMerchantByID(ctx, req.MerchantID)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanUsecase][CreateLoan] while get merchant by ID, Err: %+v", err))
		return response, err
	}
	if merchant.ID == 0 {
		return response, errors.New("merchant not found")
	}

	consumerLimit, err := uc.consumerLimitRepo.GetLimitByTenureAndConsumerID(ctx, req.Tenure, req.ConsumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanUsecase][CreateLoan] while get consumer limit by tenure and consumer ID, Err: %+v", err))
		return response, err
	}
	if consumerLimit.ID == 0 {
		return response, errors.New("consumer limit not found")
	}

	loans, err := uc.loanRepo.GetLoanByConsumerID(ctx, req.ConsumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanUsecase][CreateLoan] while get loan by consumer ID, Err: %+v", err))
		return response, err
	}

	if len(loans) > 0 {
		var (
			totalLoanAmount     float64
			totalPaidLoanAmount float64
		)

		for _, loan := range loans {
			if loan.LoanStatus != "finish" && loan.ConsumerLimitID == consumerLimit.ID {
				totalLoanAmount += loan.LoanAmount
				totalPaidLoanAmount += loan.PaidLoanAmount
			} else {
				continue
			}
		}

		remainingLimit := consumerLimit.LimitAmount - (totalLoanAmount - totalPaidLoanAmount)
		if req.LoanAmount > remainingLimit {
			errMsg := fmt.Sprintf("remaining limit: %.2f", remainingLimit)
			return response, errors.New(errMsg)
		}
	}

	dueDate := time.Now().AddDate(0, int(req.Tenure), 0)
	contractNumber := fmt.Sprintf("%d-%s-%d", req.ConsumerID, utils.GenerateUniqueString(10), req.MerchantID)
	interestAmount := req.LoanAmount * req.InterestRate / 100
	loanStatus := "on_going"

	loan := repository.Loan{
		ConsumerID:      req.ConsumerID,
		MerchantID:      req.MerchantID,
		ConsumerLimitID: consumerLimit.ID,
		LoanAmount:      req.LoanAmount,
		InterestRate:    req.InterestRate,
		InterestAmount:  interestAmount,
		LoanStatus:      loanStatus,
		DueDate:         dueDate,
		ContractNumber:  contractNumber,
		AssetName:       req.AssetName,
	}

	loanID, err := uc.loanRepo.CreateLoan(ctx, loan)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanUsecase][CreateLoan] while create loan, Err: %+v", err))
		return response, err
	}

	response = LoanResponse{
		ID:              loanID,
		ConsumerID:      req.ConsumerID,
		MerchantID:      req.MerchantID,
		ConsumerLimitID: consumerLimit.ID,
		LoanAmount:      req.LoanAmount,
		ContractNumber:  contractNumber,
		InterestRate:    req.InterestRate,
		InterestAmount:  interestAmount,
		LoanStatus:      loanStatus,
		DueDate:         dueDate.Format("2006-01-02"),
		AssetName:       req.AssetName,
	}

	return response, nil
}

func (uc *loanUsecase) GetLoanByID(ctx context.Context, loanID int64) (response LoanResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	loan, err := uc.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return response, err
	}
	if loan.ID == 0 {
		return response, errors.New("loan not found")
	}

	response = LoanResponse{
		ID:              loan.ID,
		ConsumerID:      loan.ConsumerID,
		MerchantID:      loan.MerchantID,
		ConsumerLimitID: loan.ConsumerLimitID,
		LoanAmount:      loan.LoanAmount,
		ContractNumber:  loan.ContractNumber,
		InterestRate:    loan.InterestRate,
		InterestAmount:  loan.InterestAmount,
		LoanStatus:      loan.LoanStatus,
		DueDate:         loan.DueDate.Format("2006-01-02"),
		Installment:     loan.Installment,
		AssetName:       loan.AssetName,
		CreatedAt:       loan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       loan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (uc *loanUsecase) GetLoanByConsumerID(ctx context.Context, consumerID int64) (response []LoanResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	loans, err := uc.loanRepo.GetLoanByConsumerID(ctx, consumerID)
	if err != nil {
		return response, err
	}

	for _, loan := range loans {
		response = append(response, LoanResponse{
			ID:              loan.ID,
			ConsumerID:      loan.ConsumerID,
			MerchantID:      loan.MerchantID,
			ConsumerLimitID: loan.ConsumerLimitID,
			LoanAmount:      loan.LoanAmount,
			ContractNumber:  loan.ContractNumber,
			InterestRate:    loan.InterestRate,
			InterestAmount:  loan.InterestAmount,
			LoanStatus:      loan.LoanStatus,
			DueDate:         loan.DueDate.Format("2006-01-02"),
			Installment:     loan.Installment,
			AssetName:       loan.AssetName,
			CreatedAt:       loan.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       loan.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}

func (uc *loanUsecase) DeleteLoanByID(ctx context.Context, loanID int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	loan, err := uc.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return err
	}
	if loan.ID == 0 {
		return errors.New("loan not found")
	}

	err = uc.loanRepo.DeleteLoan(ctx, loanID)
	if err != nil {
		return err
	}

	return nil
}
