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

type LoanHandler struct {
	LoanUC usecase.LoanUsecase
}

// NewLoanHandler will initialize the loan resources endpoint
func NewLoanHandler(g *echo.Group, loanUC usecase.LoanUsecase) {
	handler := &LoanHandler{
		LoanUC: loanUC,
	}

	loanGroup := g.Group("/loans")

	loanGroup.POST("", handler.Create)
	loanGroup.GET("/:id", handler.GetByID)
	loanGroup.GET("/consumer/:consumerId", handler.GetByConsumerID)
	loanGroup.DELETE("/:id", handler.Delete)
}

func (h *LoanHandler) Create(c echo.Context) error {
	req := usecase.CreateLoanRequest{}
	if err := c.Bind(&req); err != nil {
		logger.Error(fmt.Sprintf("[LoanHandler][Create] while bind request, Err: %+v", err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ConsumerID, validation.Required),
		validation.Field(&req.MerchantID, validation.Required),
		validation.Field(&req.Tenure, validation.Required),
		validation.Field(&req.LoanAmount, validation.Required),
		validation.Field(&req.InterestRate, validation.Required),
		validation.Field(&req.AssetName, validation.Required),
	); err != nil {
		logger.Warning(fmt.Sprintf("[LoanHandler][Create] while validate request, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, err.Error())
	}

	_, err := h.LoanUC.CreateLoan(c.Request().Context(), req)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusCreated, "Loan created successfully")
}

func (h *LoanHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanHandler][GetByID] while parse loan ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid loan ID")
	}

	data, err := h.LoanUC.GetLoanByID(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}

func (h *LoanHandler) GetByConsumerID(c echo.Context) error {
	consumerID, err := strconv.ParseInt(c.Param("consumerId"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanHandler][GetByConsumerID] while parse consumer ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid consumer ID")
	}

	data, err := h.LoanUC.GetLoanByConsumerID(c.Request().Context(), consumerID)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithData(c, http.StatusOK, data)
}

func (h *LoanHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("[LoanHandler][Delete] while parse loan ID, Err: %+v", err))
		return response.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid loan ID")
	}

	err = h.LoanUC.DeleteLoanByID(c.Request().Context(), id)
	if err != nil {
		return response.ErrorResponseWithMessage(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponseWithMessage(c, http.StatusOK, "Loan deleted successfully")
}
