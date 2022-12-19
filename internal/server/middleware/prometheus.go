package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/AlexeyKluev/user-balance/internal/metrics"
)

func NewPrometheusMiddleware(prometheusCollector *metrics.Collector) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wraprw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			h.ServeHTTP(wraprw, r)

			urlPath := r.URL.Path

			prometheusCollector.Vectors.ServerRequest.WithLabelValues(
				strconv.Itoa(wraprw.Status()),
				r.Method,
				urlPath,
			).Inc()

			prometheusCollector.Vectors.ServerLatency.WithLabelValues(
				strconv.Itoa(wraprw.Status()),
				r.Method,
				urlPath,
			).Observe(metrics.DurationMS(time.Since(start)))
		}

		return http.HandlerFunc(fn)
	}
}
