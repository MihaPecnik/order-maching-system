package handler

import (
	"encoding/json"
	"github.com/MihaPecnik/order-maching-system/models"
	"github.com/gorilla/mux"
	"net/http"
)

type Database interface {
	GetBottomBuy(ticker string) (models.UpdateOrderBookResponse, error)
	GetTopSell(ticker string) (models.UpdateOrderBookResponse, error)
	UpdateOrdersBook(request models.UpdateOrderBookRequest) ([]models.UpdateOrderBookResponse, error)
}

type Handler struct {
	DB Database
}

func NewHandler(database Database) *Handler {
	return &Handler{
		DB: database,
	}
}

func (h *Handler) UpdateOrderBook(w http.ResponseWriter, r *http.Request) {
	var request models.UpdateOrderBookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.DB.UpdateOrdersBook(request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)

}

func (h *Handler) GetTickerInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticker := params["ticker"]

	buy, err := h.DB.GetBottomBuy(ticker)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sell, err := h.DB.GetTopSell(ticker)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.GetTickerInfoResponse{
		Buy:  buy,
		Sell: sell,
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
