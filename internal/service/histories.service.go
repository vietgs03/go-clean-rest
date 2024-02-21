package service

import (
	model "go-test/internal/models"
	repository "go-test/internal/repository"
)

type HistoriesService interface {
	GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error)
}

type historiesServiceImpl struct {
	repo repository.HistoriesRepository
}

func NewHistoriesService(repo repository.HistoriesRepository) HistoriesService {
	return &historiesServiceImpl{repo: repo}
}

func (s *historiesServiceImpl) GetHistories(symbol, startDate, endDate, period string) ([]model.GetHistoriesResponse, error) {
	return s.repo.GetHistories(symbol, startDate, endDate, period)
}
