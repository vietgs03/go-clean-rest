package handler_test

import (
	"encoding/json"
	"go-test/internal/handler"
	model "go-test/internal/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockHistoriesService struct{}

func (m *MockHistoriesService) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
	// Return a fixed response
	return []model.GetHistoriesResponse{
		{
			Time:  time.Now(),
			Open:  1.0,
			High:  1.1,
			Low:   0.9,
			Close: 1.0,
		},
	}, nil
}

func TestGetHistoriesHandler(t *testing.T) {
	h := &handler.Handler{
		Service: &MockHistoriesService{},
	}

	req, err := http.NewRequest("GET", "/histories?start_date=2022-01-01&end_date=2022-01-02&period=1&symbol=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetHistoriesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Read the response body
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the response body into a slice of models.GetHistoriesResponse
	var response []model.GetHistoriesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the values of the response
	// Replace this with your own checks
	if len(response) != 1 {
		t.Errorf("expected 1 history, got %v", len(response))
	}
}
