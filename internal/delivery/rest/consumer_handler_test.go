package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks" // Import the mocks package
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.ConsumerUsecase) // Use the mock from the mocks package
	handler := &ConsumerHandler{ConsumerUC: mockUC}

	t.Run("successful creation", func(t *testing.T) {
		reqBody := usecase.ConsumerRequest{
			FullName:     "John Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			NIK:          "1234567890",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/consumers", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUC.On("CreateConsumer", c.Request().Context(), reqBody).Return(int64(1), nil)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Consumer created successfully")
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/consumers", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := usecase.ConsumerRequest{
			FullName:     "",
			PlaceOfBirth: "",
			DOB:          "",
			NIK:          "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/consumers", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUpdate(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.ConsumerUsecase) // Use the mock from the mocks package
	handler := &ConsumerHandler{ConsumerUC: mockUC}

	t.Run("successful update", func(t *testing.T) {
		reqBody := usecase.ConsumerRequest{
			FullName:     "John Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			NIK:          "1234567890",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/consumers/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUC.On("UpdateConsumer", c.Request().Context(), int64(1), reqBody).Return(nil)

		err := handler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Consumer updated successfully")
	})

	t.Run("invalid consumer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/consumers/invalid", bytes.NewBuffer([]byte("{}")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid consumer ID")
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/consumers/1", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := usecase.ConsumerRequest{
			FullName:     "",
			PlaceOfBirth: "",
			DOB:          "",
			NIK:          "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/consumers/1", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.ConsumerUsecase) // Use the mock from the mocks package
	handler := &ConsumerHandler{ConsumerUC: mockUC}

	t.Run("successful deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/consumers/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockUC.On("DeleteConsumer", c.Request().Context(), int64(1)).Return(nil)

		err := handler.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Consumer deleted successfully")
	})

	t.Run("invalid consumer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/consumers/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid consumer ID")
	})
}

func TestFetch(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.ConsumerUsecase) // Use the mock from the mocks package
	handler := &ConsumerHandler{ConsumerUC: mockUC}

	t.Run("successful fetch", func(t *testing.T) {
		reqBody := usecase.FetchConsumerRequest{
			// Add necessary fields for the request
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodGet, "/consumers", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockConsumers := []usecase.GetConsumerResponse{
			{ID: 1, FullName: "John Doe"},
			{ID: 2, FullName: "Jane Doe"},
		}
		mockUC.On("FetchConsumer", c.Request().Context(), reqBody).Return(mockConsumers, nil)

		err := handler.Fetch(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "John Doe")
		assert.Contains(t, rec.Body.String(), "Jane Doe")
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/consumers", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Fetch(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetByID(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.ConsumerUsecase) // Use the mock from the mocks package
	handler := &ConsumerHandler{ConsumerUC: mockUC}

	t.Run("successful get by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/consumers/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockConsumer := usecase.GetConsumerResponse{
			ID:           1,
			FullName:     "John Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			NIK:          "1234567890",
		}
		mockUC.On("GetConsumerByID", c.Request().Context(), int64(1)).Return(mockConsumer, nil)

		err := handler.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "John Doe")
	})

	t.Run("invalid consumer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/consumers/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid consumer ID")
	})
}
