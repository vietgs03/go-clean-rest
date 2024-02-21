package repository

import (
	"encoding/json"
	"fmt"
	model "go-test/internal/models"
	"io/ioutil"
	"net/http"
	"time"
)

type HistoriesRepository struct{}

func NewHistoriesRepository() *HistoriesRepository {
	return &HistoriesRepository{}
}

func (r *HistoriesRepository) processResponse(response [][]interface{}, start, end time.Time) []model.GetHistoriesResponse {
	var histories []model.GetHistoriesResponse
	for _, data := range response {
		timeStamp := int64(data[0].(float64))
		time := time.Unix(0, timeStamp*int64(time.Millisecond))

		// If the time is within the start and end dates, add it to the list
		if (time.After(start) || time.Equal(start)) && (time.Before(end) || time.Equal(end)) {
			open := data[1].(float64)
			high := data[2].(float64)
			low := data[3].(float64)
			close := data[4].(float64)

			history := model.GetHistoriesResponse{
				Time:  time,
				Open:  open,
				High:  high,
				Low:   low,
				Close: close,
			}
			histories = append(histories, history)
		}
	}
	return histories
}

func (r *HistoriesRepository) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
	// Parse the dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Convert to Unix timestamp
	startUnix := start.Unix()
	endUnix := end.Unix()

	// Create the URL
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/ohlc?vs_currency=usd&from=%d&to=%d", symbol, startUnix, endUnix)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// JSON decode
	var response [][]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Process the response
	histories := r.processResponse(response, start, end)

	return histories, nil
}
