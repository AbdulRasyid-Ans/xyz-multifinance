package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/delivery/rest"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	e := echo.New()
	mockMerchantUC := new(mocks.MerchantUsecase)
	handler := &rest.MerchantHandler{
		MerchantUC: mockMerchantUC,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockMerchantUC.On("CreateMerchant", c.Request().Context(), reqBody).Return(nil).Once()

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "Merchant created successfully")
		}
	})

	t.Run("bad request - invalid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewBuffer([]byte(`invalid`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("bad request - validation error", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "",
			MerchantType: "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockMerchantUC.On("CreateMerchant", c.Request().Context(), reqBody).Return(errors.New("internal error")).Once()

		err := handler.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	e := echo.New()
	mockMerchantUC := new(mocks.MerchantUsecase)
	handler := &rest.MerchantHandler{
		MerchantUC: mockMerchantUC,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "Updated Merchant",
			MerchantType: "Retail",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/merchants/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockMerchantUC.On("UpdateMerchant", c.Request().Context(), int64(1), reqBody).Return(nil).Once()

		err := handler.Update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Merchant updated successfully")
		}
	})

	t.Run("bad request - invalid merchant ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/merchants/invalid", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid merchant ID")
		}
	})

	t.Run("bad request - invalid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/merchants/1", bytes.NewBuffer([]byte(`invalid`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("bad request - validation error", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "",
			MerchantType: "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/merchants/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		reqBody := usecase.MerchantRequest{
			MerchantName: "Updated Merchant",
			MerchantType: "Retail",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/merchants/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockMerchantUC.On("UpdateMerchant", c.Request().Context(), int64(1), reqBody).Return(errors.New("internal error")).Once()

		err := handler.Update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestDelete(t *testing.T) {
	e := echo.New()
	mockMerchantUC := new(mocks.MerchantUsecase)
	handler := &rest.MerchantHandler{
		MerchantUC: mockMerchantUC,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/merchants/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockMerchantUC.On("DeleteMerchant", c.Request().Context(), int64(1)).Return(nil).Once()

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Merchant deleted successfully")
		}
	})

	t.Run("bad request - invalid merchant ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/merchants/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid merchant ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/merchants/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockMerchantUC.On("DeleteMerchant", c.Request().Context(), int64(1)).Return(errors.New("internal error")).Once()

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestFetch(t *testing.T) {
	e := echo.New()
	mockMerchantUC := new(mocks.MerchantUsecase)
	handler := &rest.MerchantHandler{
		MerchantUC: mockMerchantUC,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := usecase.FetchMerchantRequest{
			// Add necessary fields for FetchMerchantRequest
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodGet, "/merchants", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockData := []usecase.GetMerchantResponse{
			{
				ID:           1,
				MerchantName: "Test Merchant",
				MerchantType: "Retail",
			},
		}
		mockMerchantUC.On("FetchMerchant", c.Request().Context(), reqBody).Return(mockData, nil).Once()

		err := handler.Fetch(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Test Merchant")
		}
	})

	t.Run("bad request - invalid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants", bytes.NewBuffer([]byte(`invalid`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Fetch(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		reqBody := usecase.FetchMerchantRequest{
			// Add necessary fields for FetchMerchantRequest
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodGet, "/merchants", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockMerchantUC.On("FetchMerchant", c.Request().Context(), reqBody).Return(nil, errors.New("internal error")).Once()

		err := handler.Fetch(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	e := echo.New()
	mockMerchantUC := new(mocks.MerchantUsecase)
	handler := &rest.MerchantHandler{
		MerchantUC: mockMerchantUC,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockData := usecase.GetMerchantResponse{
			ID:           1,
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
		}
		mockMerchantUC.On("GetMerchantByID", c.Request().Context(), int64(1)).Return(mockData, nil).Once()

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Test Merchant")
		}
	})

	t.Run("bad request - invalid merchant ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid merchant ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockMerchantUC.On("GetMerchantByID", c.Request().Context(), int64(1)).Return(usecase.GetMerchantResponse{}, errors.New("internal error")).Once()

		err := handler.GetByID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
