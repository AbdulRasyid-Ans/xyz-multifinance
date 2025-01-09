package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)
	mockMerchantRepo := new(mocks.MerchantRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)

	uc := usecase.NewLoanUsecase(mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, mockMerchantRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		req := usecase.CreateLoanRequest{
			ConsumerID:   1,
			MerchantID:   1,
			Tenure:       12,
			LoanAmount:   1000,
			InterestRate: 5,
			AssetName:    "Car",
		}

		mockConsumerRepo.On("GetConsumerByID", mock.Anything, req.ConsumerID).Return(repository.Consumer{ID: 1}, nil).Once()
		mockMerchantRepo.On("GetMerchantByID", mock.Anything, req.MerchantID).Return(repository.Merchant{ID: 1}, nil).Once()
		mockConsumerLimitRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, req.Tenure, req.ConsumerID).Return(repository.ConsumerLimit{ID: 1, LimitAmount: 5000}, nil).Once()
		mockLoanRepo.On("GetLoanByConsumerID", mock.Anything, req.ConsumerID).Return([]repository.Loan{}, nil).Once()
		mockLoanRepo.On("CreateLoan", mock.Anything, mock.Anything).Return(int64(1), nil).Once()

		resp, err := uc.CreateLoan(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		mockConsumerRepo.AssertExpectations(t)
		mockMerchantRepo.AssertExpectations(t)
		mockConsumerLimitRepo.AssertExpectations(t)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("consumer not found", func(t *testing.T) {
		req := usecase.CreateLoanRequest{
			ConsumerID:   1,
			MerchantID:   1,
			Tenure:       12,
			LoanAmount:   1000,
			InterestRate: 5,
			AssetName:    "Car",
		}

		mockConsumerRepo.On("GetConsumerByID", mock.Anything, req.ConsumerID).Return(repository.Consumer{}, nil).Once()

		resp, err := uc.CreateLoan(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "consumer not found", err.Error())
		assert.Equal(t, usecase.LoanResponse{}, resp)
		mockConsumerRepo.AssertExpectations(t)
	})

	t.Run("merchant not found", func(t *testing.T) {
		req := usecase.CreateLoanRequest{
			ConsumerID:   1,
			MerchantID:   1,
			Tenure:       12,
			LoanAmount:   1000,
			InterestRate: 5,
			AssetName:    "Car",
		}

		mockConsumerRepo.On("GetConsumerByID", mock.Anything, req.ConsumerID).Return(repository.Consumer{ID: 1}, nil).Once()
		mockMerchantRepo.On("GetMerchantByID", mock.Anything, req.MerchantID).Return(repository.Merchant{}, nil).Once()

		resp, err := uc.CreateLoan(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "merchant not found", err.Error())
		assert.Equal(t, usecase.LoanResponse{}, resp)
		mockConsumerRepo.AssertExpectations(t)
		mockMerchantRepo.AssertExpectations(t)
	})

	t.Run("consumer limit not found", func(t *testing.T) {
		req := usecase.CreateLoanRequest{
			ConsumerID:   1,
			MerchantID:   1,
			Tenure:       12,
			LoanAmount:   1000,
			InterestRate: 5,
			AssetName:    "Car",
		}

		mockConsumerRepo.On("GetConsumerByID", mock.Anything, req.ConsumerID).Return(repository.Consumer{ID: 1}, nil).Once()
		mockMerchantRepo.On("GetMerchantByID", mock.Anything, req.MerchantID).Return(repository.Merchant{ID: 1}, nil).Once()
		mockConsumerLimitRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, req.Tenure, req.ConsumerID).Return(repository.ConsumerLimit{}, nil).Once()

		resp, err := uc.CreateLoan(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "consumer limit not found", err.Error())
		assert.Equal(t, usecase.LoanResponse{}, resp)
		mockConsumerRepo.AssertExpectations(t)
		mockMerchantRepo.AssertExpectations(t)
		mockConsumerLimitRepo.AssertExpectations(t)
	})
}

func TestGetLoanByID(t *testing.T) {
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)
	mockMerchantRepo := new(mocks.MerchantRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)

	uc := usecase.NewLoanUsecase(mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, mockMerchantRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		loanID := int64(1)
		loan := repository.Loan{
			ID:              loanID,
			ConsumerID:      1,
			MerchantID:      1,
			ConsumerLimitID: 1,
			LoanAmount:      1000,
			ContractNumber:  "123-abc-456",
			InterestRate:    5,
			InterestAmount:  50,
			LoanStatus:      "on_going",
			DueDate:         time.Now(),
			Installment:     12,
			AssetName:       "Car",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil).Once()

		resp, err := uc.GetLoanByID(context.Background(), loanID)
		assert.NoError(t, err)
		assert.Equal(t, loanID, resp.ID)
		assert.Equal(t, loan.ConsumerID, resp.ConsumerID)
		assert.Equal(t, loan.MerchantID, resp.MerchantID)
		assert.Equal(t, loan.ConsumerLimitID, resp.ConsumerLimitID)
		assert.Equal(t, loan.LoanAmount, resp.LoanAmount)
		assert.Equal(t, loan.ContractNumber, resp.ContractNumber)
		assert.Equal(t, loan.InterestRate, resp.InterestRate)
		assert.Equal(t, loan.InterestAmount, resp.InterestAmount)
		assert.Equal(t, loan.LoanStatus, resp.LoanStatus)
		assert.Equal(t, loan.DueDate.Format("2006-01-02"), resp.DueDate)
		assert.Equal(t, loan.Installment, resp.Installment)
		assert.Equal(t, loan.AssetName, resp.AssetName)
		assert.Equal(t, loan.CreatedAt.Format("2006-01-02 15:04:05"), resp.CreatedAt)
		assert.Equal(t, loan.UpdatedAt.Format("2006-01-02 15:04:05"), resp.UpdatedAt)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanID := int64(1)

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{}, nil).Once()

		resp, err := uc.GetLoanByID(context.Background(), loanID)
		assert.Error(t, err)
		assert.Equal(t, "loan not found", err.Error())
		assert.Equal(t, usecase.LoanResponse{}, resp)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("error fetching loan", func(t *testing.T) {
		loanID := int64(1)
		expectedErr := errors.New("unexpected error")

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{}, expectedErr).Once()

		resp, err := uc.GetLoanByID(context.Background(), loanID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, usecase.LoanResponse{}, resp)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestGetLoanByConsumerID(t *testing.T) {
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)
	mockMerchantRepo := new(mocks.MerchantRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)

	uc := usecase.NewLoanUsecase(mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, mockMerchantRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		consumerID := int64(1)
		loans := []repository.Loan{
			{
				ID:              1,
				ConsumerID:      consumerID,
				MerchantID:      1,
				ConsumerLimitID: 1,
				LoanAmount:      1000,
				ContractNumber:  "123-abc-456",
				InterestRate:    5,
				InterestAmount:  50,
				LoanStatus:      "on_going",
				DueDate:         time.Now(),
				Installment:     12,
				AssetName:       "Car",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}

		mockLoanRepo.On("GetLoanByConsumerID", mock.Anything, consumerID).Return(loans, nil).Once()

		resp, err := uc.GetLoanByConsumerID(context.Background(), consumerID)
		assert.NoError(t, err)
		assert.Len(t, resp, len(loans))
		assert.Equal(t, loans[0].ID, resp[0].ID)
		assert.Equal(t, loans[0].ConsumerID, resp[0].ConsumerID)
		assert.Equal(t, loans[0].MerchantID, resp[0].MerchantID)
		assert.Equal(t, loans[0].ConsumerLimitID, resp[0].ConsumerLimitID)
		assert.Equal(t, loans[0].LoanAmount, resp[0].LoanAmount)
		assert.Equal(t, loans[0].ContractNumber, resp[0].ContractNumber)
		assert.Equal(t, loans[0].InterestRate, resp[0].InterestRate)
		assert.Equal(t, loans[0].InterestAmount, resp[0].InterestAmount)
		assert.Equal(t, loans[0].LoanStatus, resp[0].LoanStatus)
		assert.Equal(t, loans[0].DueDate.Format("2006-01-02"), resp[0].DueDate)
		assert.Equal(t, loans[0].Installment, resp[0].Installment)
		assert.Equal(t, loans[0].AssetName, resp[0].AssetName)
		assert.Equal(t, loans[0].CreatedAt.Format("2006-01-02 15:04:05"), resp[0].CreatedAt)
		assert.Equal(t, loans[0].UpdatedAt.Format("2006-01-02 15:04:05"), resp[0].UpdatedAt)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("no loans found", func(t *testing.T) {
		consumerID := int64(1)

		mockLoanRepo.On("GetLoanByConsumerID", mock.Anything, consumerID).Return([]repository.Loan{}, nil).Once()

		resp, err := uc.GetLoanByConsumerID(context.Background(), consumerID)
		assert.NoError(t, err)
		assert.Empty(t, resp)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("error fetching loans", func(t *testing.T) {
		consumerID := int64(1)
		expectedErr := errors.New("unexpected error")

		mockLoanRepo.On("GetLoanByConsumerID", mock.Anything, consumerID).Return(nil, expectedErr).Once()

		resp, err := uc.GetLoanByConsumerID(context.Background(), consumerID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
		mockLoanRepo.AssertExpectations(t)
	})
}

func TestDeleteLoanByID(t *testing.T) {
	mockLoanRepo := new(mocks.LoanRepository)
	mockConsumerRepo := new(mocks.ConsumerRepository)
	mockMerchantRepo := new(mocks.MerchantRepository)
	mockConsumerLimitRepo := new(mocks.ConsumerLimitRepository)

	uc := usecase.NewLoanUsecase(mockLoanRepo, mockConsumerLimitRepo, mockConsumerRepo, mockMerchantRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		loanID := int64(1)

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{ID: loanID}, nil).Once()
		mockLoanRepo.On("DeleteLoan", mock.Anything, loanID).Return(nil).Once()

		err := uc.DeleteLoanByID(context.Background(), loanID)
		assert.NoError(t, err)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanID := int64(1)

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{}, nil).Once()

		err := uc.DeleteLoanByID(context.Background(), loanID)
		assert.Error(t, err)
		assert.Equal(t, "loan not found", err.Error())
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("error fetching loan", func(t *testing.T) {
		loanID := int64(1)
		expectedErr := errors.New("unexpected error")

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{}, expectedErr).Once()

		err := uc.DeleteLoanByID(context.Background(), loanID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockLoanRepo.AssertExpectations(t)
	})

	t.Run("error deleting loan", func(t *testing.T) {
		loanID := int64(1)
		expectedErr := errors.New("unexpected error")

		mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(repository.Loan{ID: loanID}, nil).Once()
		mockLoanRepo.On("DeleteLoan", mock.Anything, loanID).Return(expectedErr).Once()

		err := uc.DeleteLoanByID(context.Background(), loanID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockLoanRepo.AssertExpectations(t)
	})
}
