package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/porebric/green-api-test/internal/server"
	"github.com/porebric/green-api-test/internal/server/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := mux.NewRouter()

	handler := server.NewHandler()

	r.HandleFunc("/", handler.Home)
	r.HandleFunc("/getSettings", handler.GetSettingsHandler)
	r.HandleFunc("/getStateInstance", handler.GetStateInstanceHandler)
	r.HandleFunc("/sendMessage", handler.SendMessageHandler)
	r.HandleFunc("/sendFileByUrl", handler.SendFileByUrlHandler)

	r.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server started at :8099")

	if err := http.ListenAndServe(":8099", applyMiddlewares(r)); err != nil {
		log.Fatal("Error: ", err)
	}
}

func applyMiddlewares(handler http.Handler) http.Handler {
	handler = middlewares.RecoveryMiddleware(handler)
	handler = middlewares.RateLimitMiddleware(handler)
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.PrometheusMiddleware(handler)
	return handler
}
