package usecase

import (
	"context"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
)

type ConsumerUsecase interface {
	GetConsumerByID(ctx context.Context, id int64) (response GetConsumerResponse, err error)
	CreateConsumer(ctx context.Context, request ConsumerRequest) (id int64, err error)
	UpdateConsumer(ctx context.Context, id int64, request ConsumerRequest) (err error)
	DeleteConsumer(ctx context.Context, id int64) (err error)
	FetchConsumer(ctx context.Context, req FetchConsumerRequest) (response []GetConsumerResponse, err error)
}

type consumerUsecase struct {
	consumerRepo repository.ConsumerRepository
	ctxTimeout   time.Duration
}

type (
	GetConsumerResponse struct {
		ID           int64     `json:"id"`
		FullName     string    `json:"full_name"`
		LegalName    string    `json:"legal_name"`
		PlaceOfBirth string    `json:"place_of_birth"`
		DOB          string    `json:"dob"`
		Salary       float64   `json:"salary"`
		NIK          string    `json:"nik"`
		KTPImageURL  string    `json:"ktp_image_url"`
		SelfieURL    string    `json:"selfie_url"`
		CreatedAt    time.Time `json:"created_at"`
	}

	ConsumerRequest struct {
		FullName     string  `json:"full_name"`
		LegalName    string  `json:"legal_name"`
		PlaceOfBirth string  `json:"place_of_birth"`
		DOB          string  `json:"dob"`
		Salary       float64 `json:"salary"`
		NIK          string  `json:"nik"`
		KTPImageURL  string  `json:"ktp_image_url"`
		SelfieURL    string  `json:"selfie_url"`
	}

	FetchConsumerRequest struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
)

func NewConsumerUsecase(
	consumerRepo repository.ConsumerRepository,
	timeout time.Duration,
) ConsumerUsecase {
	return &consumerUsecase{
		consumerRepo: consumerRepo,
		ctxTimeout:   timeout,
	}
}

func (u *consumerUsecase) GetConsumerByID(ctx context.Context, id int64) (response GetConsumerResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	consumerData, err := u.consumerRepo.GetConsumerByID(ctx, id)
	if err != nil {
		return response, err
	}

	response = GetConsumerResponse{
		ID:           consumerData.ID,
		FullName:     consumerData.FullName,
		LegalName:    consumerData.LegalName,
		PlaceOfBirth: consumerData.PlaceOfBirth,
		DOB:          consumerData.DOB,
		Salary:       consumerData.Salary,
		NIK:          consumerData.NIK,
		KTPImageURL:  consumerData.KTPImageURL,
		SelfieURL:    consumerData.SelfieURL,
		CreatedAt:    consumerData.CreatedAt,
	}

	return response, nil
}

func (u *consumerUsecase) UpdateConsumer(ctx context.Context, id int64, request ConsumerRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	err = u.consumerRepo.UpdateConsumer(ctx, repository.Consumer{
		ID:           id,
		FullName:     request.FullName,
		LegalName:    request.LegalName,
		PlaceOfBirth: request.PlaceOfBirth,
		DOB:          request.DOB,
		Salary:       request.Salary,
		NIK:          request.NIK,
		KTPImageURL:  request.KTPImageURL,
		SelfieURL:    request.SelfieURL,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *consumerUsecase) DeleteConsumer(ctx context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	err = u.consumerRepo.DeleteConsumer(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *consumerUsecase) FetchConsumer(ctx context.Context, req FetchConsumerRequest) (response []GetConsumerResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	limit, offset := utils.ParsePagination(req.Page, req.Limit)

	consumerData, err := u.consumerRepo.FetchConsumer(ctx, repository.FetchConsumerRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return response, err
	}

	for _, consumer := range consumerData {
		response = append(response, GetConsumerResponse{
			ID:           consumer.ID,
			FullName:     consumer.FullName,
			LegalName:    consumer.LegalName,
			PlaceOfBirth: consumer.PlaceOfBirth,
			DOB:          consumer.DOB,
			Salary:       consumer.Salary,
			NIK:          consumer.NIK,
			KTPImageURL:  consumer.KTPImageURL,
			SelfieURL:    consumer.SelfieURL,
			CreatedAt:    consumer.CreatedAt,
		})
	}

	return response, nil
}

func (u *consumerUsecase) CreateConsumer(ctx context.Context, request ConsumerRequest) (id int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	id, err = u.consumerRepo.CreateConsumer(ctx, repository.Consumer{
		FullName:     request.FullName,
		LegalName:    request.LegalName,
		PlaceOfBirth: request.PlaceOfBirth,
		DOB:          request.DOB,
		Salary:       request.Salary,
		NIK:          request.NIK,
		KTPImageURL:  request.KTPImageURL,
		SelfieURL:    request.SelfieURL,
	})
	if err != nil {
		return id, err
	}

	return id, nil
}
