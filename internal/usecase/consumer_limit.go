package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
)

type ConsumerLimitUsecase interface {
	GetLimitByTenureAndConsumerID(ctx context.Context, tenure int16, consumerID int64) (response GetConsumerLimitResponse, err error)
	CreateOrUpdateConsumerLimit(ctx context.Context, request ConsumerLimitRequest) (err error)
	DeleteConsumerLimit(ctx context.Context, consumerLimitID int64) (err error)
	GetConsumerLimitByConsumerID(ctx context.Context, consumerID int64) (response []GetConsumerLimitResponse, err error)
}

type consumerLimitUsecase struct {
	consumerLimitRepo repository.ConsumerLimitRepository
	ctxTimeout        time.Duration
}

type (
	GetConsumerLimitResponse struct {
		ID          int64   `json:"id"`
		ConsumerID  int64   `json:"consumer_id"`
		Tenure      int16   `json:"tenure"`
		LimitAmount float64 `json:"limit_amount"`
	}

	ConsumerLimitRequest struct {
		ConsumerID  int64   `json:"consumer_id"`
		Tenure      int16   `json:"tenure"`
		LimitAmount float64 `json:"limit_amount"`
	}
)

func NewConsumerLimitUsecase(
	consumerLimitRepo repository.ConsumerLimitRepository,
	timeout time.Duration,
) ConsumerLimitUsecase {
	return &consumerLimitUsecase{
		consumerLimitRepo: consumerLimitRepo,
		ctxTimeout:        timeout,
	}
}

func (uc *consumerLimitUsecase) GetLimitByTenureAndConsumerID(ctx context.Context, tenure int16, consumerID int64) (response GetConsumerLimitResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	data, err := uc.consumerLimitRepo.GetLimitByTenureAndConsumerID(ctx, tenure, consumerID)
	if err != nil {
		return response, err
	}

	response = GetConsumerLimitResponse{
		ID:          data.ID,
		ConsumerID:  data.ConsumerID,
		Tenure:      data.Tenure,
		LimitAmount: data.LimitAmount,
	}

	return response, nil
}

func (uc *consumerLimitUsecase) CreateOrUpdateConsumerLimit(ctx context.Context, request ConsumerLimitRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	_, ok := repository.ValidConsumerLimitTenure[request.Tenure]
	if !ok {
		return errors.New("invalid tenure")
	}

	data, err := uc.consumerLimitRepo.GetLimitByTenureAndConsumerID(ctx, request.Tenure, request.ConsumerID)
	if err != nil {
		return err
	}

	if data.ID == 0 {
		_, err = uc.consumerLimitRepo.CreateConsumerLimit(ctx, repository.ConsumerLimit{
			ConsumerID:  request.ConsumerID,
			Tenure:      request.Tenure,
			LimitAmount: request.LimitAmount,
		})
		if err != nil {
			return err
		}
	} else {
		err = uc.consumerLimitRepo.UpdateConsumerLimit(ctx, repository.ConsumerLimit{
			ID:          data.ID,
			LimitAmount: request.LimitAmount,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *consumerLimitUsecase) DeleteConsumerLimit(ctx context.Context, consumerLimitID int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	err = uc.consumerLimitRepo.DeleteConsumerLimit(ctx, consumerLimitID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *consumerLimitUsecase) GetConsumerLimitByConsumerID(ctx context.Context, consumerID int64) (response []GetConsumerLimitResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	data, err := uc.consumerLimitRepo.GetConsumerLimitByConsumerID(ctx, consumerID)
	if err != nil {
		return response, err
	}

	for _, item := range data {
		response = append(response, GetConsumerLimitResponse{
			ID:          item.ID,
			ConsumerID:  item.ConsumerID,
			Tenure:      item.Tenure,
			LimitAmount: item.LimitAmount,
		})
	}

	return response, nil
}
