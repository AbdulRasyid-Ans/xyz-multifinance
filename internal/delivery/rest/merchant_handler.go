package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/response"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type MerchantHandler struct {
	MerchantUC usecase.MerchantUsecase
}

// NewMerchantHandler will initialize the merchant resources endpoint
func NewMerchantHandler(g *echo.Group, merchantUC usecase.MerchantUsecase) {
	handler := &MerchantHandler{
		MerchantUC: merchantUC,
	}

	merchantGroup := g.Group("/merchants")

	merchantGroup.POST("", handler.Create)
	merchantGroup.GET("/:id", handler.GetByID)
	merchantGroup.GET("", handler.Fetch)
	merchantGroup.PUT("/:id", handler.Update)
	merchantGroup.DELETE("/:id", handler.Delete)
}

func (h *MerchantHandler) Create(c echo.Context) error {
	req := usecase.MerchantRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][Create] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MerchantName, validation.Required),
		validation.Field(&req.MerchantType, validation.Required),
	); err != nil {
		logger.Warning(fmt.Sprintf("[MerchantHandler][Create] while validate request, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	err := h.MerchantUC.CreateMerchant(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusCreated, "Merchant created successfully")
}

func (h *MerchantHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][Update] while parse merchant ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid merchant ID")
	}

	req := usecase.MerchantRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][Update] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MerchantName, validation.Required),
		validation.Field(&req.MerchantType, validation.Required),
	); err != nil {
		logger.Warning(fmt.Sprintf("[MerchantHandler][Update] while validate request, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	err = h.MerchantUC.UpdateMerchant(c.Request().Context(), id, req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Merchant updated successfully")
}

func (h *MerchantHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][Delete] while parse merchant ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid merchant ID")
	}

	err = h.MerchantUC.DeleteMerchant(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Merchant deleted successfully")
}

func (h *MerchantHandler) Fetch(c echo.Context) error {
	req := usecase.FetchMerchantRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][Fetch] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	data, err := h.MerchantUC.FetchMerchant(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}

func (h *MerchantHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[MerchantHandler][GetByID] while parse merchant ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid merchant ID")
	}

	data, err := h.MerchantUC.GetMerchantByID(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}
