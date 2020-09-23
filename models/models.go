package models

// Tables
type Table struct {
	ID       int64   `json:"id"`
	UserId   int64   `json:"user_id"`
	Value    float64 `json:"value,string" sql:"type:decimal(10,2);"`
	Quantity int64   `json:"quantity"`
	Buy      bool    `json:"buy"`
	Ticker   string  `json:"ticker"`
}

// Request and responses
type UpdateOrderBookRequest struct {
	UserId   int64   `json:"user_id"`
	Value    float64 `json:"value,string"`
	Quantity int64   `json:"quantity"`
	Buy      bool    `json:"buy"`
	Ticker   string  `json:"ticker"`
}

type UpdateOrderBookResponse struct {
	Value    float64 `json:"value,string"`
	Quantity int64   `json:"quantity"`
}

type GetTickerInfoResponse struct {
	Buy  UpdateOrderBookResponse `json:"buy"`
	Sell UpdateOrderBookResponse `json:"sell"`
}
