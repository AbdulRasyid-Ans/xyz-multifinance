package rest

import (
	"fmt"
	"net/http"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/response"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	TransactionUC usecase.TransactionUsecase
}

// NewTransactionHandler will initialize the transaction resources endpoint
func NewTransactionHandler(g *echo.Group, transactionUC usecase.TransactionUsecase) {
	handler := &TransactionHandler{
		TransactionUC: transactionUC,
	}

	transactionGroup := g.Group("/transactions")

	transactionGroup.POST("", handler.Create)
	transactionGroup.GET("/remaining-payment", handler.GetRemainingPayment)
}

func (h *TransactionHandler) Create(c echo.Context) error {
	req := usecase.TransactionRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[TransactionHandler][Create] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ConsumerID, validation.Required),
		validation.Field(&req.LoanID, validation.Required),
		validation.Field(&req.TramsactionType, validation.Required),
	); err != nil {
		logger.Warning(fmt.Sprintf("[TransactionHandler][Create] while validate request, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	_, err := h.TransactionUC.CreateTransaction(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusCreated, "Transaction created successfully")
}

func (h *TransactionHandler) GetRemainingPayment(c echo.Context) error {
	req := usecase.TransactionRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[TransactionHandler][GetRemainingPayment] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ConsumerID, validation.Required),
		validation.Field(&req.LoanID, validation.Required),
	); err != nil {
		logger.Warning(fmt.Sprintf("[TransactionHandler][GetRemainingPayment] while validate request, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	data, err := h.TransactionUC.GetRemainingPayment(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}
