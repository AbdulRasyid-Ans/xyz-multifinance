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

func TestGetLimitByTenureAndConsumerID(t *testing.T) {
	mockRepo := new(mocks.ConsumerLimitRepository)
	uc := usecase.NewConsumerLimitUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{
			ID:          1,
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		}, nil).Once()

		ctx := context.Background()
		response, err := uc.GetLimitByTenureAndConsumerID(ctx, 2, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.ID)
		assert.Equal(t, int64(1), response.ConsumerID)
		assert.Equal(t, int16(2), response.Tenure)
		assert.Equal(t, 1000.0, response.LimitAmount)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{}, errors.New("some error")).Once()

		ctx := context.Background()
		response, err := uc.GetLimitByTenureAndConsumerID(ctx, 2, 1)

		assert.Error(t, err)
		assert.Equal(t, usecase.GetConsumerLimitResponse{}, response)

		mockRepo.AssertExpectations(t)
	})
}

func TestCreateOrUpdateConsumerLimit(t *testing.T) {
	mockRepo := new(mocks.ConsumerLimitRepository)
	uc := usecase.NewConsumerLimitUsecase(mockRepo, time.Second*2)

	t.Run("success create", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{}, nil).Once()
		mockRepo.On("CreateConsumerLimit", mock.Anything, mock.Anything).Return(int64(1), nil).Once()

		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success update", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{
			ID:          1,
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 500.0,
		}, nil).Once()
		mockRepo.On("UpdateConsumerLimit", mock.Anything, mock.Anything).Return(nil).Once()

		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid tenure", func(t *testing.T) {
		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      99,
			LimitAmount: 1000.0,
		})

		assert.Error(t, err)
		assert.Equal(t, "invalid tenure", err.Error())
	})

	t.Run("repository error on GetLimitByTenureAndConsumerID", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{}, errors.New("some error")).Once()

		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		})

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on CreateConsumerLimit", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{}, nil).Once()
		mockRepo.On("CreateConsumerLimit", mock.Anything, mock.Anything).Return(int64(0), errors.New("some error")).Once()

		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		})

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on UpdateConsumerLimit", func(t *testing.T) {
		mockRepo.On("GetLimitByTenureAndConsumerID", mock.Anything, int16(2), int64(1)).Return(repository.ConsumerLimit{
			ID:          1,
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 500.0,
		}, nil).Once()
		mockRepo.On("UpdateConsumerLimit", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		ctx := context.Background()
		err := uc.CreateOrUpdateConsumerLimit(ctx, usecase.ConsumerLimitRequest{
			ConsumerID:  1,
			Tenure:      2,
			LimitAmount: 1000.0,
		})

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteConsumerLimit(t *testing.T) {
	mockRepo := new(mocks.ConsumerLimitRepository)
	uc := usecase.NewConsumerLimitUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteConsumerLimit", mock.Anything, int64(1)).Return(nil).Once()

		ctx := context.Background()
		err := uc.DeleteConsumerLimit(ctx, int64(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteConsumerLimit", mock.Anything, int64(1)).Return(errors.New("some error")).Once()

		ctx := context.Background()
		err := uc.DeleteConsumerLimit(ctx, int64(1))

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetConsumerLimitByConsumerID(t *testing.T) {
	mockRepo := new(mocks.ConsumerLimitRepository)
	uc := usecase.NewConsumerLimitUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetConsumerLimitByConsumerID", mock.Anything, int64(1)).Return([]repository.ConsumerLimit{
			{
				ID:          1,
				ConsumerID:  1,
				Tenure:      2,
				LimitAmount: 1000.0,
			},
			{
				ID:          2,
				ConsumerID:  1,
				Tenure:      3,
				LimitAmount: 2000.0,
			},
		}, nil).Once()

		ctx := context.Background()
		response, err := uc.GetConsumerLimitByConsumerID(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, int64(1), response[0].ID)
		assert.Equal(t, int64(1), response[0].ConsumerID)
		assert.Equal(t, int16(2), response[0].Tenure)
		assert.Equal(t, 1000.0, response[0].LimitAmount)
		assert.Equal(t, int64(2), response[1].ID)
		assert.Equal(t, int64(1), response[1].ConsumerID)
		assert.Equal(t, int16(3), response[1].Tenure)
		assert.Equal(t, 2000.0, response[1].LimitAmount)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetConsumerLimitByConsumerID", mock.Anything, int64(1)).Return(nil, errors.New("some error")).Once()

		ctx := context.Background()
		response, err := uc.GetConsumerLimitByConsumerID(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "some error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
