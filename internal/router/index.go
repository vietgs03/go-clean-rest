package router

import (
	"go-test/internal/handler"
	"go-test/internal/repository"
	"go-test/internal/service"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()

	historiesRepository := repository.NewHistoriesRepository()
	historiesService := service.NewHistoriesService(*historiesRepository)
	historiesHandler := handler.NewHistoriesHandler(historiesService)
	router.HandleFunc("/get_histories", historiesHandler.GetHistoriesHandler) // call get_histories endpoint
	return router
}
