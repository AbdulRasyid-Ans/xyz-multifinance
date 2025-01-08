package usecase

import (
	"context"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
)

type MerchantUsecase interface {
	GetMerchantByID(ctx context.Context, id int64) (response GetMerchantResponse, err error)
	UpdateMerchant(ctx context.Context, id int64, request MerchantRequest) (err error)
	DeleteMerchant(ctx context.Context, id int64) (err error)
	FetchMerchant(ctx context.Context, req FetchMerchantRequest) (response []GetMerchantResponse, err error)
	CreateMerchant(ctx context.Context, request MerchantRequest) (err error)
}

type merchantUsecase struct {
	merchantRepo repository.MerchantRepository
	ctxTimeout   time.Duration
}

type (
	GetMerchantResponse struct {
		ID           int64     `json:"id"`
		MerchantName string    `json:"merchant_name"`
		MerchantType string    `json:"merchant_type"`
		CreatedAt    time.Time `json:"created_at"`
	}

	MerchantRequest struct {
		MerchantName string `json:"merchant_name"`
		MerchantType string `json:"merchant_type"`
	}

	FetchMerchantRequest struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
)

func NewMerchantUsecase(
	merchantRepo repository.MerchantRepository,
	timeout time.Duration,
) MerchantUsecase {
	return &merchantUsecase{
		merchantRepo: merchantRepo,
		ctxTimeout:   timeout,
	}
}

func (uc *merchantUsecase) GetMerchantByID(ctx context.Context, id int64) (response GetMerchantResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	data, err := uc.merchantRepo.GetMerchantByID(ctx, id)
	if err != nil {
		return GetMerchantResponse{}, err
	}

	response.ID = data.ID
	response.MerchantName = data.MerchantName
	response.MerchantType = data.MerchantType
	response.CreatedAt = data.CreatedAt

	return response, nil
}

func (uc *merchantUsecase) UpdateMerchant(ctx context.Context, id int64, request MerchantRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	data := repository.Merchant{
		ID:           id,
		MerchantName: request.MerchantName,
		MerchantType: request.MerchantType,
	}

	err = uc.merchantRepo.UpdateMerchant(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (uc *merchantUsecase) DeleteMerchant(ctx context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	err = uc.merchantRepo.DeleteMerchant(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *merchantUsecase) FetchMerchant(ctx context.Context, req FetchMerchantRequest) (response []GetMerchantResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	limit, offset := utils.ParsePagination(req.Page, req.Limit)

	data, err := uc.merchantRepo.FetchMerchant(ctx, repository.FetchMerchantRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return response, err
	}

	for _, item := range data {
		response = append(response, GetMerchantResponse{
			ID:           item.ID,
			MerchantName: item.MerchantName,
			MerchantType: item.MerchantType,
			CreatedAt:    item.CreatedAt,
		})
	}

	return response, nil
}

func (uc *merchantUsecase) CreateMerchant(ctx context.Context, request MerchantRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	data := repository.Merchant{
		MerchantName: request.MerchantName,
		MerchantType: request.MerchantType,
	}

	_, err = uc.merchantRepo.CreateMerchant(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
