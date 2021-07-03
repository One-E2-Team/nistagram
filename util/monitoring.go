package util

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
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

var (
	totalRequests, requestInitiators, responseStatus *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
)

func InitMonitoring(ms string)  {

	if !DockerChecker() {
		return
	}

	initMonitoringVars(ms)

	prometheus.Register(totalRequests)
	prometheus.Register(requestInitiators)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	router := mux.NewRouter().StrictSlash(true)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()

			totalRequests.WithLabelValues(path).Inc()

			rw := NewResponseWriter(w)

			start := time.Now()
			next.ServeHTTP(rw, r)
			elapsed := time.Since(start)

			statusCode := rw.statusCode

			httpDuration.WithLabelValues(path, strconv.Itoa(statusCode)).Observe(elapsed.Seconds())
			responseStatus.WithLabelValues(path, strconv.Itoa(statusCode)).Inc()

			initiator := rw.Header().Get("initiator")
			requestInitiators.WithLabelValues(path, initiator).Inc()

		})
	})

	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe(":9090", router)
}

func initMonitoringVars(ms string) {
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nistagram_" + ms + "_http_requests_total",
			Help: "Number of requests.",
		},
		[]string{"path"},
	)

	requestInitiators = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nistagram_" + ms + "_http_requests_total",
			Help: "Initiators of HTTP requests.",
		},
		[]string{"path", "initiator"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nistagram_" + ms + "_http_response_status",
			Help: "Status of HTTP response",
		},
		[]string{"path", "status"},
	)

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "nistagram_" + ms + "_http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "status"})
}