package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	StatusOk              string = "ok"
	StatusError           string = "error"
	StatusAccessDenied    string = "access_denied"
	StatusValidationError string = "validation_error"
	StatusNotFound        string = "not_found"
)

// Коллектор метрик
type Collector struct {
	Registry *prometheus.Registry
	Vectors  struct {
		ServerRequest       *prometheus.CounterVec
		ServerLatency       *prometheus.HistogramVec
		GraphqlQueryLatency *prometheus.HistogramVec
		HTTPClientLatency   *prometheus.HistogramVec
		DBQueryLatency      *prometheus.HistogramVec
		RedisCommandLatency *prometheus.HistogramVec
	}
}

func NewCollector(serviceName string) *Collector {
	collectors := Collector{
		Registry: prometheus.NewRegistry(),
	}

	collectors.registerVectors(serviceName)

	return &collectors
}

func (p *Collector) registerVectors(serviceName string) {
	var (
		dflBuckets = []float64{10, 50, 100, 300, 500, 1000, 5000, 10000, 30000}
	)

	p.Vectors.ServerRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "server_requests_total",
			Help:        "Сколько HTTP-запросов обрабатывается, разделяется по статусу, методу и HTTP-пути.",
			ConstLabels: prometheus.Labels{"service": serviceName},
		},
		[]string{"code", "method", "path"},
	)
	p.Registry.MustRegister(p.Vectors.ServerRequest)

	p.Vectors.ServerLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_request_duration_milliseconds",
		Help: "Сколько времени потребовалось для обработки запроса, " +
			"разбитого по коду статуса, методу и HTTP-пути.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     dflBuckets,
	},
		[]string{"code", "method", "path"},
	)
	p.Registry.MustRegister(p.Vectors.ServerLatency)

	p.Vectors.GraphqlQueryLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "graphql_query_duration_milliseconds",
		Help: "Сколько времени потребовалось на обработку запросов Graphql сервису, " +
			"разбитого по коду статуса, методу.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     dflBuckets,
	},
		[]string{"code", "method"},
	)
	p.Registry.MustRegister(p.Vectors.GraphqlQueryLatency)

	p.Vectors.HTTPClientLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_client_request_duration_milliseconds",
		Help: "Сколько времени потребовалось на запросы к сторонним сервисам, " +
			"разбитого по коду статуса, методу и HTTP-пути.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     dflBuckets,
	},
		[]string{"code", "method", "path"},
	)
	p.Registry.MustRegister(p.Vectors.HTTPClientLatency)

	p.Vectors.DBQueryLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "sql_query_duration_milliseconds",
		Help: "Сколько времени потребовалось для обработки SQL запросов, " +
			"разбитого по коду статуса, методу.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     dflBuckets,
	},
		[]string{"query"},
	)
	p.Registry.MustRegister(p.Vectors.DBQueryLatency)

	p.Vectors.RedisCommandLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "redis_command_duration_milliseconds",
		Help: "Сколько времени потребовалось для выполнения комманды в Redis, " +
			"разбитого по комманде и ключу.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     append([]float64{.1, .25, .5, 1, 3, 5}, dflBuckets...),
	},
		[]string{"command"},
	)
	p.Registry.MustRegister(p.Vectors.RedisCommandLatency)
}

// shortcut
func (p *Collector) GraphqlQueryLatencyObserve(start time.Time, lvs ...string) {
	p.Vectors.GraphqlQueryLatency.
		WithLabelValues(lvs...).
		Observe(DurationMS(time.Since(start)))
}

func DurationMS(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()) / 1e6
}

func Since(start time.Time) float64 {
	return DurationMS(time.Since(start))
}
