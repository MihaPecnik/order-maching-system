package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
}

func (s *Server) Serve(port string) error {
	return http.ListenAndServe(port, s.router)
}

type Handler interface {
	UpdateOrderBook(w http.ResponseWriter, r *http.Request)
	GetTickerInfo(w http.ResponseWriter, r *http.Request)
}


func NewServer(h Handler) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/orderbook", h.UpdateOrderBook).Methods("PUT")
	router.HandleFunc("/api/v1/tickerinfo/{ticker}", h.GetTickerInfo).Methods("GET")

	return &Server{
		router: router,
	}
}
