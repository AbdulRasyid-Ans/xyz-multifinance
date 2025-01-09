package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRemainingPayment(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)

	uc := usecase.NewTransactionUsecase(mockTransactionRepo, mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, time.Second*2)

	tests := []struct {
		name    string
		req     usecase.TransactionRequest
		setup   func()
		want    usecase.RemainingPaymentResponse
		wantErr bool
	}{
		{
			name: "invalid transaction type",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "invalid",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "consumer not found",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan not found",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan not belong to consumer",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{ID: 1, ConsumerID: 2}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan already finished",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{ID: 1, ConsumerID: 1, LoanStatus: "finish"}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "consumer limit not found",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{ID: 1, ConsumerID: 1, ConsumerLimitID: 1}, nil).Once()
				mockConsumerLimitRepo.On("GetConsumerLimitByID", mock.Anything, int64(1)).Return(repository.ConsumerLimit{}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "successful full payment",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "full",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{
					ID:                 1,
					ConsumerID:         1,
					ConsumerLimitID:    1,
					ContractNumber:     "123",
					LoanAmount:         1000,
					PaidLoanAmount:     500,
					InterestAmount:     100,
					PaidInterestAmount: 50,
					DueDate:            time.Now().AddDate(0, 1, 0),
				}, nil).Once()
				mockConsumerLimitRepo.On("GetConsumerLimitByID", mock.Anything, int64(1)).Return(repository.ConsumerLimit{ID: 1, Tenure: 12, ConsumerID: 1}, nil).Once()
			},
			want: usecase.RemainingPaymentResponse{
				ContractNumber:          "123",
				Tenure:                  12,
				DueDate:                 time.Now().AddDate(0, 1, 0),
				PaidLoanAmount:          500,
				PaidInterestAmount:      50,
				RemainingLoanAmount:     500,
				RemainingInterestAmount: 50,
				TotalRemainingAmount:    550,
			},
			wantErr: false,
		},
		{
			name: "successful installment payment",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{
					ID:                 1,
					ConsumerID:         1,
					ConsumerLimitID:    1,
					ContractNumber:     "123",
					LoanAmount:         1200,
					PaidLoanAmount:     100,
					InterestAmount:     120,
					PaidInterestAmount: 10,
					DueDate:            time.Now().AddDate(0, 1, 0),
					Installment:        1,
				}, nil).Once()
				mockConsumerLimitRepo.On("GetConsumerLimitByID", mock.Anything, int64(1)).Return(repository.ConsumerLimit{ID: 1, Tenure: 12, ConsumerID: 1}, nil).Once()
			},
			want: usecase.RemainingPaymentResponse{
				ContractNumber:          "123",
				Installment:             2,
				Tenure:                  12,
				DueDate:                 time.Now().AddDate(0, 1, 0),
				PaidLoanAmount:          100,
				PaidInterestAmount:      10,
				RemainingLoanAmount:     100,
				RemainingInterestAmount: 10,
				TotalRemainingAmount:    110,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := uc.GetRemainingPayment(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ContractNumber, got.ContractNumber)
				assert.Equal(t, tt.want.Installment, got.Installment)
				assert.Equal(t, tt.want.Tenure, got.Tenure)
				assert.WithinDuration(t, tt.want.DueDate, got.DueDate, time.Second)
				assert.Equal(t, tt.want.PaidLoanAmount, got.PaidLoanAmount)
				assert.Equal(t, tt.want.PaidInterestAmount, got.PaidInterestAmount)
				assert.Equal(t, tt.want.RemainingLoanAmount, got.RemainingLoanAmount)
				assert.Equal(t, tt.want.RemainingInterestAmount, got.RemainingInterestAmount)
				assert.Equal(t, tt.want.TotalRemainingAmount, got.TotalRemainingAmount)
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	mockTransactionRepo := new(mocks.TransactionRepository)
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)

	uc := usecase.NewTransactionUsecase(mockTransactionRepo, mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, time.Second*2)

	tests := []struct {
		name    string
		req     usecase.TransactionRequest
		setup   func()
		want    usecase.GetTransactionResponse
		wantErr bool
	}{
		{
			name: "invalid transaction type",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "invalid",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "consumer not found",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan not found",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan not belong to consumer",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{ID: 1, ConsumerID: 2}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "loan already finished",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{ID: 1, ConsumerID: 1, LoanStatus: "finish"}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "successful full payment",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "full",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{
					ID:                 1,
					ConsumerID:         1,
					ConsumerLimitID:    1,
					ContractNumber:     "123",
					LoanAmount:         1000,
					PaidLoanAmount:     500,
					InterestAmount:     100,
					PaidInterestAmount: 50,
					DueDate:            time.Now().AddDate(0, 1, 0),
				}, nil).Once()
				mockConsumerLimitRepo.On("GetConsumerLimitByID", mock.Anything, int64(1)).Return(repository.ConsumerLimit{ID: 1, Tenure: 12, ConsumerID: 1}, nil).Once()
				mockTransactionRepo.On("BeginTx", mock.Anything).Return(nil, nil).Once()
				mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				mockLoanRepo.On("UpdateLoan", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mockTransactionRepo.On("CommitTx", mock.Anything, mock.Anything).Return(nil).Once()
			},
			want: usecase.GetTransactionResponse{
				ID:          1,
				ConsumerID:  1,
				LoanID:      1,
				Amount:      550,
				Description: "Payment for loan 123",
			},
			wantErr: false,
		},
		{
			name: "successful installment payment",
			req: usecase.TransactionRequest{
				ConsumerID:      1,
				LoanID:          1,
				TramsactionType: "installment",
			},
			setup: func() {
				mockConsumerRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(repository.Consumer{ID: 1}, nil).Once()
				mockLoanRepo.On("GetLoanByID", mock.Anything, int64(1)).Return(repository.Loan{
					ID:                 1,
					ConsumerID:         1,
					ConsumerLimitID:    1,
					ContractNumber:     "123",
					LoanAmount:         1200,
					PaidLoanAmount:     100,
					InterestAmount:     120,
					PaidInterestAmount: 10,
					DueDate:            time.Now().AddDate(0, 1, 0),
					Installment:        1,
				}, nil).Once()
				mockConsumerLimitRepo.On("GetConsumerLimitByID", mock.Anything, int64(1)).Return(repository.ConsumerLimit{ID: 1, Tenure: 12, ConsumerID: 1}, nil).Once()
				mockTransactionRepo.On("BeginTx", mock.Anything).Return(nil, nil).Once()
				mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				mockLoanRepo.On("UpdateLoan", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mockTransactionRepo.On("CommitTx", mock.Anything, mock.Anything).Return(nil).Once()
			},
			want: usecase.GetTransactionResponse{
				ID:          1,
				ConsumerID:  1,
				LoanID:      1,
				Amount:      110,
				Description: "Payment for loan 123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := uc.CreateTransaction(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.ConsumerID, got.ConsumerID)
				assert.Equal(t, tt.want.LoanID, got.LoanID)
				assert.Equal(t, tt.want.Amount, got.Amount)
				assert.Equal(t, tt.want.Description, got.Description)
			}
		})
	}
}
