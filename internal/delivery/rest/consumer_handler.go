package rest

import (
	"net/http"
	"strconv"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/response"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type ConsumerHandler struct {
	ConsumerUC usecase.ConsumerUsecase
}

// NewConsumerHandler will initialize the consumer resources endpoint
func NewConsumerHandler(g *echo.Group, consumerUC usecase.ConsumerUsecase) {
	handler := &ConsumerHandler{
		ConsumerUC: consumerUC,
	}

	consumerGroup := g.Group("/consumers")

	consumerGroup.POST("", handler.Create)
	consumerGroup.GET("/:id", handler.GetByID)
	consumerGroup.GET("", handler.Fetch)
	consumerGroup.PUT("/:id", handler.Update)
	consumerGroup.DELETE("/:id", handler.Delete)
}

func (h *ConsumerHandler) Create(c echo.Context) error {
	req := usecase.ConsumerRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.FullName, validation.Required),
		validation.Field(&req.PlaceOfBirth, validation.Required),
		validation.Field(&req.DOB, validation.Required),
		validation.Field(&req.NIK, validation.Required),
	); err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	_, err := h.ConsumerUC.CreateConsumer(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusCreated, "Consumer created successfully")
}

func (h *ConsumerHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	req := usecase.ConsumerRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.FullName, validation.Required),
		validation.Field(&req.PlaceOfBirth, validation.Required),
		validation.Field(&req.DOB, validation.Required),
		validation.Field(&req.NIK, validation.Required),
	); err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	err = h.ConsumerUC.UpdateConsumer(c.Request().Context(), id, req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Consumer updated successfully")
}

func (h *ConsumerHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	err = h.ConsumerUC.DeleteConsumer(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Consumer deleted successfully")
}

func (h *ConsumerHandler) Fetch(c echo.Context) error {
	req := usecase.FetchConsumerRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	consumers, err := h.ConsumerUC.FetchConsumer(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, consumers)
}

func (h *ConsumerHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	consumer, err := h.ConsumerUC.GetConsumerByID(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	if consumer.ID == 0 {
		return response.ErrorResponseWithMessage(c, http.StatusNotFound, "Consumer not found")
	}

	return response.SuccessResponseWithData(c, http.StatusOK, consumer)
}
