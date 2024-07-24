package main

import (
	"fmt"
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

	fmt.Println("Server started at :8080")

	http.ListenAndServe(":8090", applyMiddlewares(r))
}

func applyMiddlewares(handler http.Handler) http.Handler {
	handler = middlewares.RecoveryMiddleware(handler)
	handler = middlewares.RateLimitMiddleware(handler)
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.PrometheusMiddleware(handler)
	return handler
}
