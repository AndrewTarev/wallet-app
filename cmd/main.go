package main

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	logger "github.com/sirupsen/logrus"

	"wallet-app/internal/app/delivery/http"
	"wallet-app/internal/app/repository"
	"wallet-app/internal/app/services"
	"wallet-app/internal/configs"
	"wallet-app/internal/infrastructure/database"
	logging "wallet-app/internal/infrastructure/logger"
	"wallet-app/internal/infrastructure/server"
)

// @title           Wallet
// @version         1.0
// @description     API для операций с кошельком.

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./internal/configs")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	// Настройка логгера
	logging.SetupLogger(&cfg.Logging)

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(cfg.Database.Dsn)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close()

	applyMigrations(cfg.Database.Dsn)

	repo := repository.NewRepository(dbConn)
	service := services.NewService(repo)
	handlers := http.NewHandler(service)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRoutes())
}

func applyMigrations(dsn string) {
	m, err := migrate.New(
		"file:///app/internal/infrastructure/database/migrations",
		dsn,
	)
	if err != nil {
		logger.Fatalf("Could not initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatalf("Could not apply migrations: %v", err)
	}

	logger.Debug("Applied migrations")
}
