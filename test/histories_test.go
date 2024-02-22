package handler_test

import (
	"bytes"
	model "go-test/internal/models"
	repository "go-test/internal/repository"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type MockHistoriesService struct {
	GetHistoriesCallCount int
}

func (m *MockHistoriesService) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
	m.GetHistoriesCallCount++
	return []model.GetHistoriesResponse{
		{
			Time:   time.Now().Unix(),
			Open:   1.0,
			High:   1.1,
			Low:    0.9,
			Close:  1.0,
			Change: 0.0,
		},
	}, nil
}

func TestGetHistoriesResponseFormat(t *testing.T) {

	repository := repository.NewHistoriesRepository()

	// Execute the API call
	histories, err := repository.GetHistories("bitcoin", "2024-02-15T00:00:00Z", "2024-02-20T00:00:00Z", "30M")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Assertions
	if _, ok := interface{}(histories[0].High).(float64); !ok {
		t.Error("Expected 'high' field to be of type float64")
	}
	if _, ok := interface{}(histories[0].Low).(float64); !ok {
		t.Error("Expected 'low' field to be of type float64")
	}
	if _, ok := interface{}(histories[0].Open).(float64); !ok {
		t.Error("Expected 'open' field to be of type float64")
	}
	if _, ok := interface{}(histories[0].Close).(float64); !ok {
		t.Error("Expected 'close' field to be of type float64")
	}
	if _, ok := interface{}(histories[0].Time).(int64); !ok {
		t.Error("Expected 'time' field to be of type int64")
	}
	if _, ok := interface{}(histories[0].Change).(float64); !ok {
		t.Error("Expected 'change' field to be of type float64")
	}
}

func TestCacheUsage(t *testing.T) {

	repository := repository.NewHistoriesRepository()

	// Execute the API call for the first time
	_, err := repository.GetHistories("bitcoin", "2024-02-15T00:00:00Z", "2024-02-20T00:00:00Z", "30M")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// ensure api not call again
	httpClient := &http.Client{Transport: &mockTransport{}}
	repository.SetHttpClient(httpClient)

	// Execute the API call for the second time
	_, err = repository.GetHistories("bitcoin", "2024-02-15T00:00:00Z", "2024-02-20T00:00:00Z", "30M")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
}

type mockTransport struct {
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Return an empty response without making a real HTTP request
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}, nil
}
