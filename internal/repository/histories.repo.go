package repository

import (
	"context"
	"encoding/json"
	"fmt"
	model "go-test/internal/models"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type HistoriesRepository struct{}

func NewHistoriesRepository() *HistoriesRepository {
	return &HistoriesRepository{}
}

var ctx = context.Background()

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r *HistoriesRepository) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/ohlc?vs_currency=usd&days=7&precision=18", symbol)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	layout := "2006-01-02T15:04:05Z"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, err
	}
	// JSON decode
	var response [][]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	// get interval
	interval := getIntervalbyPeriod(period)
	var lastIncluded time.Time
	var histories []model.GetHistoriesResponse
	for _, data := range response {
		timeStamp := int64(data[0].(float64))
		open := data[1].(float64)
		high := data[2].(float64)
		low := data[3].(float64)
		close := data[4].(float64)

		// Convert timestamp to time.Time
		time := time.Unix(0, timeStamp*int64(time.Millisecond))
		if time.Before(start) || time.After(end) || (!lastIncluded.IsZero() && time.Sub(lastIncluded) < interval) {
			continue
		}
		// Add the history to the list
		history := model.GetHistoriesResponse{
			Time:  timeStamp,
			Open:  open,
			High:  high,
			Low:   low,
			Close: close,
		}
		histories = append(histories, history)
	}

	return histories, nil
}

func getIntervalbyPeriod(period string) time.Duration {
	var interval time.Duration
	switch period {
	case "30M":
		interval = 30 * time.Minute
	case "1H":
		interval = 1 * time.Hour
	case "4H":
		interval = 4 * time.Hour
	case "1D":
		interval = 24 * time.Hour
	default:
		interval = 4 * time.Hour
	}
	return interval
}
