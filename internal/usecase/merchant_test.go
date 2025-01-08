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

func TestGetMerchantByID(t *testing.T) {
	mockRepo := new(mocks.MerchantRepository)
	uc := usecase.NewMerchantUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockMerchant := repository.Merchant{
			ID:           1,
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
			CreatedAt:    time.Now(),
		}

		mockRepo.On("GetMerchantByID", mock.Anything, int64(1)).Return(mockMerchant, nil)

		ctx := context.Background()
		response, err := uc.GetMerchantByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, mockMerchant.ID, response.ID)
		assert.Equal(t, mockMerchant.MerchantName, response.MerchantName)
		assert.Equal(t, mockMerchant.MerchantType, response.MerchantType)
		assert.WithinDuration(t, mockMerchant.CreatedAt, response.CreatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetMerchantByID", mock.Anything, int64(2)).Return(repository.Merchant{}, errors.New("merchant not found"))

		ctx := context.Background()
		response, err := uc.GetMerchantByID(ctx, 2)

		assert.Error(t, err)
		assert.Equal(t, usecase.GetMerchantResponse{}, response)

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateMerchant(t *testing.T) {
	mockRepo := new(mocks.MerchantRepository)
	uc := usecase.NewMerchantUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRequest := usecase.MerchantRequest{
			MerchantName: "Updated Merchant",
			MerchantType: "Wholesale",
		}

		mockRepo.On("UpdateMerchant", mock.Anything, repository.Merchant{
			ID:           1,
			MerchantName: mockRequest.MerchantName,
			MerchantType: mockRequest.MerchantType,
		}).Return(nil)

		ctx := context.Background()
		err := uc.UpdateMerchant(ctx, 1, mockRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRequest := usecase.MerchantRequest{
			MerchantName: "Updated Merchant",
			MerchantType: "Wholesale",
		}

		mockRepo.On("UpdateMerchant", mock.Anything, repository.Merchant{
			ID:           2,
			MerchantName: mockRequest.MerchantName,
			MerchantType: mockRequest.MerchantType,
		}).Return(errors.New("update failed"))

		ctx := context.Background()
		err := uc.UpdateMerchant(ctx, 2, mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteMerchant(t *testing.T) {
	mockRepo := new(mocks.MerchantRepository)
	uc := usecase.NewMerchantUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteMerchant", mock.Anything, int64(1)).Return(nil)

		ctx := context.Background()
		err := uc.DeleteMerchant(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteMerchant", mock.Anything, int64(2)).Return(errors.New("delete failed"))

		ctx := context.Background()
		err := uc.DeleteMerchant(ctx, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestFetchMerchant(t *testing.T) {
	mockRepo := new(mocks.MerchantRepository)
	uc := usecase.NewMerchantUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockMerchants := []repository.Merchant{
			{
				ID:           1,
				MerchantName: "Merchant 1",
				MerchantType: "Retail",
				CreatedAt:    time.Now(),
			},
			{
				ID:           2,
				MerchantName: "Merchant 2",
				MerchantType: "Wholesale",
				CreatedAt:    time.Now(),
			},
		}

		mockRepo.On("FetchMerchant", mock.Anything, repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}).Return(mockMerchants, nil).Once()

		ctx := context.Background()
		req := usecase.FetchMerchantRequest{
			Page:  1,
			Limit: 10,
		}
		response, err := uc.FetchMerchant(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, mockMerchants[0].ID, response[0].ID)
		assert.Equal(t, mockMerchants[0].MerchantName, response[0].MerchantName)
		assert.Equal(t, mockMerchants[0].MerchantType, response[0].MerchantType)
		assert.WithinDuration(t, mockMerchants[0].CreatedAt, response[0].CreatedAt, time.Second)
		assert.Equal(t, mockMerchants[1].ID, response[1].ID)
		assert.Equal(t, mockMerchants[1].MerchantName, response[1].MerchantName)
		assert.Equal(t, mockMerchants[1].MerchantType, response[1].MerchantType)
		assert.WithinDuration(t, mockMerchants[1].CreatedAt, response[1].CreatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("FetchMerchant", mock.Anything, repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}).Return(nil, errors.New("fetch failed")).Once()

		ctx := context.Background()
		req := usecase.FetchMerchantRequest{
			Page:  1,
			Limit: 10,
		}
		response, err := uc.FetchMerchant(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)

		mockRepo.AssertExpectations(t)
	})
}

func TestCreateMerchant(t *testing.T) {
	mockRepo := new(mocks.MerchantRepository)
	uc := usecase.NewMerchantUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRequest := usecase.MerchantRequest{
			MerchantName: "New Merchant",
			MerchantType: "Retail",
		}

		mockRepo.On("CreateMerchant", mock.Anything, repository.Merchant{
			MerchantName: mockRequest.MerchantName,
			MerchantType: mockRequest.MerchantType,
		}).Return(int64(1), nil).Once()

		ctx := context.Background()
		err := uc.CreateMerchant(ctx, mockRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRequest := usecase.MerchantRequest{
			MerchantName: "New Merchant",
			MerchantType: "Retail",
		}

		mockRepo.On("CreateMerchant", mock.Anything, repository.Merchant{
			MerchantName: mockRequest.MerchantName,
			MerchantType: mockRequest.MerchantType,
		}).Return(int64(0), errors.New("create failed")).Once()

		ctx := context.Background()
		err := uc.CreateMerchant(ctx, mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
