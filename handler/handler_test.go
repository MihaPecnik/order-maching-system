package handler

import (
	"encoding/json"
	"github.com/MihaPecnik/order-maching-system/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type dbMock struct {
	mock.Mock
}

func (d *dbMock) GetBottomBuy(ticker string) (models.UpdateOrderBookResponse, error) {
	args := d.Called(ticker)
	return args.Get(0).(models.UpdateOrderBookResponse), args.Error(1)
}

func (d *dbMock) GetTopSell(ticker string) (models.UpdateOrderBookResponse, error) {
	args := d.Called(ticker)
	return args.Get(0).(models.UpdateOrderBookResponse), args.Error(1)
}

func (d *dbMock) UpdateOrdersBook(request models.UpdateOrderBookRequest) ([]models.UpdateOrderBookResponse, error) {
	args := d.Called(request)
	return args.Get(0).([]models.UpdateOrderBookResponse), args.Error(1)
}

func TestNewHandler(t *testing.T) {
	dbMock := &dbMock{}
	handler := NewHandler(dbMock)
	assert.NotNil(t, handler)
}

func TestGetTickerInfoSuccess(t *testing.T) {
	dbMock := &dbMock{}
	order1 := models.UpdateOrderBookResponse{
		Value:    20.0,
		Quantity: 15,
	}

	order2 := models.UpdateOrderBookResponse{
		Value:    19.0,
		Quantity: 10,
	}

	dbMock.On("GetBottomBuy", mock.Anything).Return(order1, nil)
	dbMock.On("GetTopSell", mock.Anything).Return(order2, nil)

	h := NewHandler(dbMock)
	req, err := http.NewRequest("Get", "/api/v1/tickerinfo/ticker", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	vars := map[string]string{
		"ticker": "APPL",
	}

	req = mux.SetURLVars(req, vars)
	h.GetTickerInfo(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusOK)

	expected := `{"buy":{"value":"20","quantity":15},"sell":{"value":"19","quantity":10}}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetTickerInfoFailBuy(t *testing.T) {
	dbMock := &dbMock{}
	dbMock.On("GetBottomBuy", mock.Anything).
		Return(models.UpdateOrderBookResponse{}, errors.New("internal error"))

	h := NewHandler(dbMock)
	req, err := http.NewRequest("Get", "/api/v1/tickerinfo/ticker", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	vars := map[string]string{
		"ticker": "APPL",
	}

	req = mux.SetURLVars(req, vars)
	h.GetTickerInfo(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusInternalServerError)

	expected := `{"error":"internal error"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetTickerInfoFailSell(t *testing.T) {
	dbMock := &dbMock{}
	order1 := models.UpdateOrderBookResponse{
		Value:    20.0,
		Quantity: 15,
	}

	dbMock.On("GetBottomBuy", mock.Anything).Return(order1, nil)
	dbMock.On("GetTopSell", mock.Anything).
		Return(models.UpdateOrderBookResponse{}, errors.New("internal error"))

	h := NewHandler(dbMock)
	req, err := http.NewRequest("Get", "/api/v1/tickerinfo/ticker", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	vars := map[string]string{
		"ticker": "APPL",
	}

	req = mux.SetURLVars(req, vars)
	h.GetTickerInfo(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusInternalServerError)

	expected := `{"error":"internal error"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestUpdateOrderBookSuccess(t *testing.T) {
	dbMock := &dbMock{}
	response := []models.UpdateOrderBookResponse{
		{
			Value:    200.1,
			Quantity: 2,
		},
		{
			Value:    200,
			Quantity: 2,
		},
	}

	request := models.UpdateOrderBookRequest{
		UserId:   1,
		Value:    199,
		Quantity: 7,
		Buy:      false,
		Ticker:   "APPL",
	}

	dbMock.On("UpdateOrdersBook", request).Return(response, nil)

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(request)
	req, err := http.NewRequest("PUT", "/api/v1/orderbook", strings.NewReader(string(bytes)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.UpdateOrderBook(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusOK)

	expected := `[{"value":"200.1","quantity":2},{"value":"200","quantity":2}]`
	assert.Equal(t, expected, rr.Body.String())
}

func TestUpdateOrderBookFail(t *testing.T) {
	dbMock := &dbMock{}
	response := []models.UpdateOrderBookResponse{
		{
			Value:    200.1,
			Quantity: 2,
		},
		{
			Value:    200,
			Quantity: 2,
		},
	}

	request := models.UpdateOrderBookRequest{
		UserId:   1,
		Value:    199,
		Quantity: 7,
		Buy:      false,
		Ticker:   "APPL",
	}

	dbMock.On("UpdateOrdersBook", request).Return(response, nil)

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(request)
	req, err := http.NewRequest("PUT", "/api/v1/orderbook", strings.NewReader(string(bytes)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.UpdateOrderBook(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusOK)

	expected := `[{"value":"200.1","quantity":2},{"value":"200","quantity":2}]`
	assert.Equal(t, expected, rr.Body.String())
}

func TestArchiveNoteFail(t *testing.T) {
	dbMock := &dbMock{}
	request := models.UpdateOrderBookRequest{
		UserId:   1,
		Value:    199,
		Quantity: 7,
		Buy:      false,
		Ticker:   "APPL",
	}
	bytes, _ := json.Marshal(request)

	dbMock.On("UpdateOrdersBook", request).Return([]models.UpdateOrderBookResponse{}, errors.New("error"))

	req, err := http.NewRequest("PUT", "/api/v1/orderbook", strings.NewReader(string("")))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h := NewHandler(dbMock)
	h.UpdateOrderBook(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)

	expected := `{"error":"EOF"}`
	assert.Equal(t, expected, rr.Body.String())

	// ////

	req1, err := http.NewRequest("PUT", "/api/v1/orderbook", strings.NewReader(string(bytes)))
	require.NoError(t, err)
	rr1 := httptest.NewRecorder()
	h1 := NewHandler(dbMock)

	h1.UpdateOrderBook(rr1, req1)

	assert.Equal(t, rr1.Code, http.StatusInternalServerError)
	expected = `{"error":"error"}`
	assert.Equal(t, expected, rr1.Body.String())
}
