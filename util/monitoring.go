package util

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
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

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error){
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack error")
	}else{
		return h.Hijack()
	}
}

var (
	totalRequests, requestInitiators, responseStatus *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
)

func InitMonitoring(ms string, router *mux.Router)  {

	if !DockerChecker() {
		return
	}

	initMonitoringVars(ms)

	prometheus.Register(totalRequests)
	prometheus.Register(requestInitiators)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	promRouter := mux.NewRouter().StrictSlash(true)

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

	promRouter.Path("/metrics").Handler(promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9090", promRouter)
		if err != nil {
			fmt.Println("FAILED TO START PROMETHEUS METRICS!")
			return
		}
	}()
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
			Name: "nistagram_http_requests_initiators",
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