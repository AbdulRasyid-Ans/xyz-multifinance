package rest

import (
	"net/http"
	"strconv"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/response"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type ConsumerLimitHandler struct {
	ConsumerLimitUC usecase.ConsumerLimitUsecase
}

// NewConsumerLimitHandler will initialize the consumer limit resources endpoint
func NewConsumerLimitHandler(g *echo.Group, consumerLimitUC usecase.ConsumerLimitUsecase) {
	handler := &ConsumerLimitHandler{
		ConsumerLimitUC: consumerLimitUC,
	}

	consumerLimitGroup := g.Group("/consumer-limits")

	consumerLimitGroup.GET("/:consumerId", handler.GetByConsumerID)
	consumerLimitGroup.GET("/:consumerId/:tenure", handler.GetByConsumerIDAndTenure)
	consumerLimitGroup.POST("", handler.CreateOrUpdate)
	consumerLimitGroup.DELETE("/:id", handler.Delete)
}

func (h *ConsumerLimitHandler) GetByConsumerID(c echo.Context) error {
	consumerID, err := strconv.ParseInt(c.Param("consumerId"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	data, err := h.ConsumerLimitUC.GetConsumerLimitByConsumerID(c.Request().Context(), consumerID)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}

func (h *ConsumerLimitHandler) GetByConsumerIDAndTenure(c echo.Context) error {
	consumerID, err := strconv.ParseInt(c.Param("consumerId"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	tenure, err := strconv.ParseInt(c.Param("tenure"), 10, 16)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid tenure")
	}

	data, err := h.ConsumerLimitUC.GetLimitByTenureAndConsumerID(c.Request().Context(), int16(tenure), consumerID)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}

func (h *ConsumerLimitHandler) CreateOrUpdate(c echo.Context) error {
	req := usecase.ConsumerLimitRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ConsumerID, validation.Required),
		validation.Field(&req.Tenure, validation.Required),
		validation.Field(&req.LimitAmount, validation.Required),
	); err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	err := h.ConsumerLimitUC.CreateOrUpdateConsumerLimit(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Consumer limit created/updated successfully")
}

func (h *ConsumerLimitHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer limit ID")
	}

	err = h.ConsumerLimitUC.DeleteConsumerLimit(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Consumer limit deleted successfully")
}
