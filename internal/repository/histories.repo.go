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

type HistoriesRepository struct {
	redisClient *redis.Client
	httpClient  *http.Client
}

func NewHistoriesRepository() *HistoriesRepository {
	return &HistoriesRepository{
		redisClient: NewClient(),
		httpClient:  &http.Client{},
	}
}

var ctx = context.Background()

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
func (r *HistoriesRepository) SetHttpClient(client *http.Client) {
	r.httpClient = client
}
func (r *HistoriesRepository) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
	// Tạo key cho cache Redis
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", symbol, startDate, endDate, period)

	// Kiểm tra cache trong Redis
	cachedHistories, err := r.GetCachedHistories(cacheKey)
	if err == nil {
		return cachedHistories, nil
	}

	// Gọi API Coingecko nếu không có cache hoặc cache đã hết hạn
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/ohlc?vs_currency=usd&days=7&precision=18", symbol)
	res, err := r.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Xử lý dữ liệu từ API
	histories, err := r.ProcessAPIResponse(body, startDate, endDate, period)
	if err != nil {
		return nil, err
	}

	// Lưu cache vào Redis
	err = r.SetCachedHistories(cacheKey, histories)
	if err != nil {
		// Xử lý lỗi khi lưu cache
		return nil, err
	}

	return histories, nil
}

func (r *HistoriesRepository) GetCachedHistories(cacheKey string) ([]model.GetHistoriesResponse, error) {
	// Lấy dữ liệu từ Redis
	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	// Decode dữ liệu từ cache
	var cachedHistories []model.GetHistoriesResponse
	if err := json.Unmarshal([]byte(cachedData), &cachedHistories); err != nil {
		return nil, err
	}

	return cachedHistories, nil
}

func (r *HistoriesRepository) SetCachedHistories(cacheKey string, histories []model.GetHistoriesResponse) error {
	// Encode - save redis
	jsonData, err := json.Marshal(histories)
	if err != nil {
		return err
	}

	// time out 1 h
	expiration := time.Hour
	if err := r.redisClient.Set(ctx, cacheKey, jsonData, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *HistoriesRepository) ProcessAPIResponse(body []byte, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
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
	var lastClose float64
	var histories []model.GetHistoriesResponse
	for _, data := range response {
		timeStamp := int64(data[0].(float64))
		open := data[1].(float64)
		high := data[2].(float64)
		low := data[3].(float64)
		close := data[4].(float64)

		var change float64
		if lastClose != 0 {
			change = (close - lastClose) / lastClose * 100
		}
		lastClose = close

		// Convert timestamp to time.Time
		time := time.Unix(0, timeStamp*int64(time.Millisecond))
		if time.Before(start) || time.After(end) || (!lastIncluded.IsZero() && time.Sub(lastIncluded) < interval) {
			continue
		}
		// Add the history to the list
		history := model.GetHistoriesResponse{
			Time:   timeStamp,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Change: change,
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
