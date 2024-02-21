package model

import "time"

// History struct
type History struct {
	ID     int       `json:"id"`
	Symbol string    `json:"symbol"`
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	Hight  float64   `json:"hight"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
}

type GetHistoriesRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Period    string    `json:"period"`
	Symbol    string    `json:"symbol"`
}

type GetHistoriesResponse struct {
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	Time   time.Time `json:"time"`
	Change float64   `json:"change"`
}
