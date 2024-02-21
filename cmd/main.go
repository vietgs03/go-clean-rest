package main

import (
	"go-test/internal/router"
	"net/http"
)

func main() {
	r := router.SetupRouter()
	http.ListenAndServe(":8080", r)
}
