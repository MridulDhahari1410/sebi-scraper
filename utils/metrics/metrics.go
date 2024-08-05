package metrics

import (
	"strconv"

	"sebi-scrapper/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// BucketConfig is used to initialise bucket.
type BucketConfig struct {
	Start float64
	Width float64
	Count int
}

// Config is used to configure the metrics.
type Config struct {
	Bucket BucketConfig
}

// timers.
var (
	dbQueryTimer             *prometheus.HistogramVec
	httpRequestTimer         *prometheus.HistogramVec
	externalHTTPRequestTimer *prometheus.HistogramVec
)

// counters.
var (
	httpRequestCounter          *prometheus.CounterVec
	httpTotalRequestCounter     *prometheus.CounterVec
	httpResponseStatusCounter   *prometheus.CounterVec
	externalHTTPResponseCounter *prometheus.CounterVec
)

// Init is used to initialise metrics.
func Init(cfg Config) {
	httpRequestTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "httpResponseTimeSeconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.LinearBuckets(cfg.Bucket.Start, cfg.Bucket.Width, cfg.Bucket.Count),
	}, []string{"path"})
	externalHTTPRequestTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "externalHttpRequestTimer",
		Help:    "Duration of External HTTP requests.",
		Buckets: prometheus.LinearBuckets(cfg.Bucket.Start, cfg.Bucket.Width, cfg.Bucket.Count),
	}, []string{"name"})
	dbQueryTimer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "DAOTime",
		Help:    "DAO time",
		Buckets: prometheus.LinearBuckets(cfg.Bucket.Start, cfg.Bucket.Width, cfg.Bucket.Count),
	}, []string{"model", "query", "method"})

	httpTotalRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "httpRequestsTotal",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)
	httpRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routerRequestsTotal",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"status", "host", "path", "method"},
	)
	httpResponseStatusCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "responseStatus",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)
	externalHTTPResponseCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "externalCount",
			Help: "How many external HTTP requests are processed.",
		},
		[]string{"name"},
	)
}

// GetMetricsMiddleware is to add prometheus timer and counter stats for requests.
func GetMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler := utils.GetHandlerName(ctx)
		timer := prometheus.NewTimer(httpRequestTimer.WithLabelValues(handler))
		httpTotalRequestCounter.WithLabelValues(handler).Inc()
		ctx.Next()
		httpResponseStatusCounter.WithLabelValues(strconv.Itoa(ctx.Writer.Status())).Inc()
		httpRequestCounter.WithLabelValues(strconv.Itoa(ctx.Writer.Status()), ctx.Request.Host, handler,
			ctx.Request.Method).Inc()
		timer.ObserveDuration()
	}
}

// GetDBQueryTimer is to log the time taken to run a database request.
func GetDBQueryTimer(name string, model string, query string) *prometheus.Timer {
	return prometheus.NewTimer(dbQueryTimer.WithLabelValues(model, query, name))
}

// GetExternalHTTPRequestTimer is to log the time taken to run a database request.
func GetExternalHTTPRequestTimer(name string) *prometheus.Timer {
	return prometheus.NewTimer(externalHTTPRequestTimer.WithLabelValues(name))
}

// GetExternalHTTPResponseCounter is to log the time taken to run a database request.
func GetExternalHTTPResponseCounter(name string) *prometheus.Counter {
	externalHTTPResponseCounter.WithLabelValues(name).Inc()
	return nil
}

// HTTPMetrics is the wrapper to add metrics to HTTP requests.
func HTTPMetrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// IncrementExternalHTTPResponseCounter is to log the number of external requests.
func IncrementExternalHTTPResponseCounter(name string) {
	externalHTTPResponseCounter.WithLabelValues(name).Inc()
}
