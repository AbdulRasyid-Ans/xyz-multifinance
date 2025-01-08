package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/delivery/rest"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/mocks"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByConsumerID(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.ConsumerLimitUsecase)
	handler := &rest.ConsumerLimitHandler{
		ConsumerLimitUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		consumerID := int64(1)
		expectedData := []usecase.GetConsumerLimitResponse{
			{
				ID:          1,
				ConsumerID:  10,
				Tenure:      2,
				LimitAmount: 1000000,
			},
		}
		mockUsecase.On("GetConsumerLimitByConsumerID", mock.Anything, consumerID).Return(expectedData, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/"+strconv.FormatInt(consumerID, 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId")
		c.SetParamValues(strconv.FormatInt(consumerID, 10))

		err := handler.GetByConsumerID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "success")
		}
	})

	t.Run("invalid consumer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/invalid", nil)
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
		consumerID := int64(1)
		mockUsecase.On("GetConsumerLimitByConsumerID", mock.Anything, consumerID).Return(nil, errors.New("some error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/"+strconv.FormatInt(consumerID, 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId")
		c.SetParamValues(strconv.FormatInt(consumerID, 10))

		err := handler.GetByConsumerID(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}

func TestGetByConsumerIDAndTenure(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.ConsumerLimitUsecase)
	handler := &rest.ConsumerLimitHandler{
		ConsumerLimitUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		consumerID := int64(1)
		tenure := int16(2)
		expectedData := usecase.GetConsumerLimitResponse{
			ID:          1,
			ConsumerID:  10,
			Tenure:      2,
			LimitAmount: 1000000,
		}
		mockUsecase.On("GetLimitByTenureAndConsumerID", mock.Anything, tenure, consumerID).Return(expectedData, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/"+strconv.FormatInt(consumerID, 10)+"/"+strconv.FormatInt(int64(tenure), 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId", "tenure")
		c.SetParamValues(strconv.FormatInt(consumerID, 10), strconv.FormatInt(int64(tenure), 10))

		err := handler.GetByConsumerIDAndTenure(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "success")
		}
	})

	t.Run("invalid consumer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/invalid/2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId", "tenure")
		c.SetParamValues("invalid", "2")

		err := handler.GetByConsumerIDAndTenure(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid consumer ID")
		}
	})

	t.Run("invalid tenure", func(t *testing.T) {
		consumerID := int64(1)
		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/"+strconv.FormatInt(consumerID, 10)+"/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId", "tenure")
		c.SetParamValues(strconv.FormatInt(consumerID, 10), "invalid")

		err := handler.GetByConsumerIDAndTenure(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid tenure")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		consumerID := int64(1)
		tenure := int16(2)
		mockUsecase.On("GetLimitByTenureAndConsumerID", mock.Anything, tenure, consumerID).Return(usecase.GetConsumerLimitResponse{}, errors.New("some error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/consumer-limits/"+strconv.FormatInt(consumerID, 10)+"/"+strconv.FormatInt(int64(tenure), 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("consumerId", "tenure")
		c.SetParamValues(strconv.FormatInt(consumerID, 10), strconv.FormatInt(int64(tenure), 10))

		err := handler.GetByConsumerIDAndTenure(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}

func TestCreateOrUpdateConsumerLimit(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.ConsumerLimitUsecase)
	handler := &rest.ConsumerLimitHandler{
		ConsumerLimitUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "tenure": 2, "limit_amount": 1000000}`
		req := httptest.NewRequest(http.MethodPost, "/consumer-limits", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("CreateOrUpdateConsumerLimit", mock.Anything, mock.Anything).Return(nil).Once()

		err := handler.CreateOrUpdate(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Consumer limit created/updated successfully")
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		reqBody := `invalid body`
		req := httptest.NewRequest(http.MethodPost, "/consumer-limits", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateOrUpdate(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "invalid character")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1}`
		req := httptest.NewRequest(http.MethodPost, "/consumer-limits", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateOrUpdate(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "cannot be blank")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		reqBody := `{"consumer_id": 1, "tenure": 2, "limit_amount": 1000000}`
		req := httptest.NewRequest(http.MethodPost, "/consumer-limits", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("CreateOrUpdateConsumerLimit", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		err := handler.CreateOrUpdate(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}
func TestDeleteConsumeLimit(t *testing.T) {
	e := echo.New()
	mockUsecase := new(mocks.ConsumerLimitUsecase)
	handler := &rest.ConsumerLimitHandler{
		ConsumerLimitUC: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		id := int64(1)
		mockUsecase.On("DeleteConsumerLimit", mock.Anything, id).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/consumer-limits/"+strconv.FormatInt(id, 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Consumer limit deleted successfully")
		}
	})

	t.Run("invalid consumer limit ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/consumer-limits/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid consumer limit ID")
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		id := int64(1)
		mockUsecase.On("DeleteConsumerLimit", mock.Anything, id).Return(errors.New("some error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/consumer-limits/"+strconv.FormatInt(id, 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		err := handler.Delete(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "some error")
		}
	})
}
