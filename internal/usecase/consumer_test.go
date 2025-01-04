package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	uc "github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetConsumerByID(t *testing.T) {
	mockRepo := new(mocks.ConsumerRepository)
	usecase := uc.NewConsumerUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockConsumer := repository.Consumer{
			ID:           1,
			FullName:     "John Doe",
			LegalName:    "Johnathan Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			Salary:       50000,
			NIK:          "1234567890",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
			CreatedAt:    time.Now(),
		}

		mockRepo.On("GetConsumerByID", mock.Anything, int64(1)).Return(mockConsumer, nil)

		ctx := context.Background()
		response, err := usecase.GetConsumerByID(ctx, int64(1))

		assert.NoError(t, err)
		assert.Equal(t, mockConsumer.ID, response.ID)
		assert.Equal(t, mockConsumer.FullName, response.FullName)
		assert.Equal(t, mockConsumer.LegalName, response.LegalName)
		assert.Equal(t, mockConsumer.PlaceOfBirth, response.PlaceOfBirth)
		assert.Equal(t, mockConsumer.DOB, response.DOB)
		assert.Equal(t, mockConsumer.Salary, response.Salary)
		assert.Equal(t, mockConsumer.NIK, response.NIK)
		assert.Equal(t, mockConsumer.KTPImageURL, response.KTPImageURL)
		assert.Equal(t, mockConsumer.SelfieURL, response.SelfieURL)
		assert.WithinDuration(t, mockConsumer.CreatedAt, response.CreatedAt, time.Second)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetConsumerByID", mock.Anything, int64(2)).Return(repository.Consumer{}, errors.New("consumer not found"))

		ctx := context.Background()
		response, err := usecase.GetConsumerByID(ctx, int64(2))

		assert.Error(t, err)
		assert.Equal(t, int64(0), response.ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateConsumer(t *testing.T) {
	mockRepo := new(mocks.ConsumerRepository)
	usecase := uc.NewConsumerUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockConsumerRequest := uc.ConsumerRequest{
			FullName:     "Jane Doe",
			LegalName:    "Jane Doe",
			PlaceOfBirth: "Los Angeles",
			DOB:          "1992-02-02",
			Salary:       60000,
			NIK:          "0987654321",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
		}

		mockRepo.On("CreateConsumer", mock.Anything, mock.AnythingOfType("repository.Consumer")).Return(int64(1), nil)

		ctx := context.Background()
		id, err := usecase.CreateConsumer(ctx, mockConsumerRequest)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		mockRepo.AssertExpectations(t)
	})
}

func TestFetchConsumer(t *testing.T) {
	mockRepo := new(mocks.ConsumerRepository)
	usecase := uc.NewConsumerUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockConsumers := []repository.Consumer{
			{
				ID:           1,
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "New York",
				DOB:          "1990-01-01",
				Salary:       50000,
				NIK:          "1234567890",
				KTPImageURL:  "http://example.com/ktp.jpg",
				SelfieURL:    "http://example.com/selfie.jpg",
				CreatedAt:    time.Now(),
			},
			{
				ID:           2,
				FullName:     "Jane Doe",
				LegalName:    "Jane Doe",
				PlaceOfBirth: "Los Angeles",
				DOB:          "1992-02-02",
				Salary:       60000,
				NIK:          "0987654321",
				KTPImageURL:  "http://example.com/ktp.jpg",
				SelfieURL:    "http://example.com/selfie.jpg",
				CreatedAt:    time.Now(),
			},
		}

		mockRepo.On("FetchConsumer", mock.Anything, mock.AnythingOfType("repository.FetchConsumerRequest")).Return(mockConsumers, nil)

		ctx := context.Background()
		req := uc.FetchConsumerRequest{
			Page:  1,
			Limit: 10,
		}
		response, err := usecase.FetchConsumer(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, mockConsumers[0].ID, response[0].ID)
		assert.Equal(t, mockConsumers[0].FullName, response[0].FullName)
		assert.Equal(t, mockConsumers[0].LegalName, response[0].LegalName)
		assert.Equal(t, mockConsumers[0].PlaceOfBirth, response[0].PlaceOfBirth)
		assert.Equal(t, mockConsumers[0].DOB, response[0].DOB)
		assert.Equal(t, mockConsumers[0].Salary, response[0].Salary)
		assert.Equal(t, mockConsumers[0].NIK, response[0].NIK)
		assert.Equal(t, mockConsumers[0].KTPImageURL, response[0].KTPImageURL)
		assert.Equal(t, mockConsumers[0].SelfieURL, response[0].SelfieURL)
		assert.WithinDuration(t, mockConsumers[0].CreatedAt, response[0].CreatedAt, time.Second)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateConsumer(t *testing.T) {
	mockRepo := new(mocks.ConsumerRepository)
	usecase := uc.NewConsumerUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockConsumerRequest := uc.ConsumerRequest{
			FullName:     "Jane Doe",
			LegalName:    "Jane Doe",
			PlaceOfBirth: "Los Angeles",
			DOB:          "1992-02-02",
			Salary:       60000,
			NIK:          "0987654321",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
		}

		mockRepo.On("UpdateConsumer", mock.Anything, mock.AnythingOfType("repository.Consumer")).Return(nil)

		ctx := context.Background()
		err := usecase.UpdateConsumer(ctx, int64(1), mockConsumerRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteConsumer(t *testing.T) {
	mockRepo := new(mocks.ConsumerRepository)
	usecase := uc.NewConsumerUsecase(mockRepo, time.Second*2)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteConsumer", mock.Anything, int64(1)).Return(nil)

		ctx := context.Background()
		err := usecase.DeleteConsumer(ctx, int64(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteConsumer", mock.Anything, int64(2)).Return(errors.New("consumer not found"))

		ctx := context.Background()
		err := usecase.DeleteConsumer(ctx, int64(2))

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
