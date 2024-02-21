package handler

import (
	"encoding/json"
	service "go-test/internal/service"
	"net/http"
)

type Handler struct {
	Service service.HistoriesService
}

func NewHistoriesHandler(s service.HistoriesService) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) GetHistoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Get param
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")
	period := r.FormValue("period")
	symbol := r.FormValue("symbol")

	// Call service
	histories, err := h.Service.GetHistories(symbol, startDate, endDate, period)
	if err != nil {
		http.Error(w, "Failed to get histories", http.StatusInternalServerError)
		return
	}

	// Return json data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(histories)
}
