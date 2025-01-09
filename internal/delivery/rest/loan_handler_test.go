package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/delivery/rest"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.LoanUsecase)
	handler := &rest.LoanHandler{
		LoanUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := `{"consumer_id":1,"merchant_id":1,"tenure":12,"loan_amount":1000000,"interest_rate":5.5,"asset_name":"Car"}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("CreateLoan", mock.Anything, mock.AnythingOfType("usecase.CreateLoanRequest")).Return(usecase.LoanResponse{}, nil).Once()

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "Loan created successfully")
		}
	})

	t.Run("bind error", func(t *testing.T) {
		reqBody := `invalid body`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "invalid character")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"consumer_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "cannot be blank")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		reqBody := `{"consumer_id":1,"merchant_id":1,"tenure":12,"loan_amount":1000000,"interest_rate":5.5,"asset_name":"Car"}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("CreateLoan", mock.Anything, mock.AnythingOfType("usecase.CreateLoanRequest")).Return(usecase.LoanResponse{}, errors.New("some error")).Once()

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}

func TestGetLoanByID(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.LoanUsecase)
	handler := &rest.LoanHandler{
		LoanUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUsecase.On("GetLoanByID", mock.Anything, int64(1)).Return(usecase.LoanResponse{}, nil).Once()

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), `"data":`)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid loan ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUsecase.On("GetLoanByID", mock.Anything, int64(1)).Return(usecase.LoanResponse{}, errors.New("some error")).Once()

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}

func TestGetLoanByConsumerID(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.LoanUsecase)
	handler := &rest.LoanHandler{
		LoanUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/consumer/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId")
		c.SetParamValues("1")

		mockUsecase.On("GetLoanByConsumerID", mock.Anything, int64(1)).Return([]usecase.LoanResponse{}, nil).Once()

		err := handler.GetByConsumerID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), `"data":`)
		}
	})

	t.Run("invalid consumer id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/consumer/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId")
		c.SetParamValues("invalid")

		err := handler.GetByConsumerID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid consumer ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/consumer/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId")
		c.SetParamValues("1")

		mockUsecase.On("GetLoanByConsumerID", mock.Anything, int64(1)).Return([]usecase.LoanResponse{}, errors.New("some error")).Once()

		err := handler.GetByConsumerID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}

func TestDeleteLoan(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.LoanUsecase)
	handler := &rest.LoanHandler{
		LoanUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/loans/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUsecase.On("DeleteLoanByID", mock.Anything, int64(1)).Return(nil).Once()

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Loan deleted successfully")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/loans/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid loan ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/loans/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUsecase.On("DeleteLoanByID", mock.Anything, int64(1)).Return(errors.New("some error")).Once()

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}
