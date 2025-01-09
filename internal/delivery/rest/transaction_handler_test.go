package rest_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/delivery/rest"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	e := echo.New()
	mockTransactionUC := new(mocks.TransactionUsecase)
	handler := &rest.TransactionHandler{
		TransactionUC: mockTransactionUC,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "loan_id": 1, "transaction_type": "payment"}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTransactionUC.On("CreateTransaction", mock.Anything, mock.Anything).Return(usecase.GetTransactionResponse{}, nil).Once()

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Transaction created successfully")
	})

	t.Run("bind error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString("invalid body"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "loan_id": 1, "transaction_type": "payment"}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTransactionUC.On("CreateTransaction", mock.Anything, mock.Anything).Return(usecase.GetTransactionResponse{}, errors.New("usecase error")).Once()

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetRemainingPayment(t *testing.T) {
	e := echo.New()
	mockTransactionUC := new(mocks.TransactionUsecase)
	handler := &rest.TransactionHandler{
		TransactionUC: mockTransactionUC,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "loan_id": 1}`
		req := httptest.NewRequest(http.MethodGet, "/transactions/remaining-payment", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTransactionUC.On("GetRemainingPayment", mock.Anything, mock.Anything).Return(usecase.RemainingPaymentResponse{}, nil).Once()

		err := handler.GetRemainingPayment(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions/remaining-payment", bytes.NewBufferString("invalid body"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetRemainingPayment(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1}`
		req := httptest.NewRequest(http.MethodGet, "/transactions/remaining-payment", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetRemainingPayment(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "loan_id": 1}`
		req := httptest.NewRequest(http.MethodGet, "/transactions/remaining-payment", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTransactionUC.On("GetRemainingPayment", mock.Anything, mock.Anything).Return(usecase.RemainingPaymentResponse{}, errors.New("usecase error")).Once()

		err := handler.GetRemainingPayment(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
