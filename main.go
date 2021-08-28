package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of GET requests",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	},
	[]string{"path"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequest.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(totalRequest)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func main() {
	router := mux.NewRouter()

	router.Use(prometheusMiddleware)
	
	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	
	log.Println("Serving request on port 9000: localhost:9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatalln(err)
}