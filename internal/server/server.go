package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // swagger for chi
	"go.uber.org/zap"
	"gopkg.in/tylerb/graceful.v1"

	"github.com/AlexeyKluev/user-balance/internal/app"
	"github.com/AlexeyKluev/user-balance/internal/server/handlers"
	middlewares "github.com/AlexeyKluev/user-balance/internal/server/middleware"
)

const CORSMaxAgeSeconds = 300

type Server struct {
	router *chi.Mux
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		router: chi.NewRouter(),
		logger: logger,
	}
}

func (s *Server) InitMiddlewares(resources *app.Resources) {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(6))

	options := cors.Options{
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Time-Zone",
		},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           CORSMaxAgeSeconds,
	}

	s.router.Use(middleware.Logger)

	options.AllowedOrigins = []string{"*"} // TODO: add origins

	cors := cors.New(options)
	s.router.Use(cors.Handler)

	s.router.Use(middlewares.NewPrometheusMiddleware(resources.MetricsCollector))
}

func (s *Server) InitRoutes(resources *app.Resources) {
	s.router.Get("/probes/readiness", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("i'm ready")); err != nil {
			resources.Logger.Error("failed to write response", zap.Error(err))
		}
	})

	s.router.Get("/probes/liveness", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("i'm alive")); err != nil {
			resources.Logger.Error("failed to write response", zap.Error(err))
		}
	})

	s.router.Handle("/metrics", promhttp.HandlerFor(resources.MetricsCollector.Registry, promhttp.HandlerOpts{}))

	s.router.Get("/users/{id:[0-9]+}/balance", handlers.NewUserBalanceHandler(resources))
	s.router.Post("/users/{id:[0-9]+}/accural", handlers.NewAccrualFundsHandler(resources))
	s.router.Post("/users/{id:[0-9]+}/reservation", handlers.NewReservationFundsHandler(resources))

	// /swagger/index.html
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
	))
}

func (s *Server) ListenAndServe(addr string, shutdownInitiated func()) error {
	s.logger.Info("ListenAndServe", zap.String("addr", addr))

	const shutdownTimeout = 10 * time.Second

	srv := &graceful.Server{
		Timeout:           shutdownTimeout,
		ShutdownInitiated: shutdownInitiated,
		Server: &http.Server{
			Addr:    addr,
			Handler: s.router,
		},
	}

	return srv.ListenAndServe()
}
