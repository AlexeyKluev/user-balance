package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/app"
	"github.com/AlexeyKluev/user-balance/internal/config"
	"github.com/AlexeyKluev/user-balance/internal/server"
	"github.com/AlexeyKluev/user-balance/internal/version"
)

// @title           Swagger User-balance service API
// @version         0.0.1

// @contact.name   Aleksey Klyuev
// @contact.email  welcome@adklyuev.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func main() {
	// Загружаем env-переменные
	_ = godotenv.Load()

	// Создаем конфиг
	cfg, err := config.InitConfig(version.App)
	if err != nil {
		log.Fatal(err)
	}

	//// Создаем ресурсы
	resources, err := app.NewResources(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем веб-сервер
	srv := server.NewServer(resources.Logger)
	srv.InitMiddlewares(resources)
	srv.InitRoutes(resources)

	_, cancelCtx := context.WithCancel(context.Background())

	beforeShutdown := func() {
		resources.Logger.Info("Остановка сервера...")

		cancelCtx()

		if !resources.Config.IsProduction {
			const shutdownIdle = 5 * time.Second

			time.Sleep(shutdownIdle)
		}

		if err := resources.Close(); err != nil {
			resources.Logger.Error("Не удалось закрыть ресурсы", zap.Error(err))
		}
	}

	// Запускаем веб-сервер
	if err = srv.ListenAndServe(resources.Config.Addr, beforeShutdown); err != nil {
		resources.Logger.Fatal("Ошибка при запуске сервера", zap.Error(err))
	}
}
