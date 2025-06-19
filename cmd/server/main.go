package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"news-service/internal/cache"
	"news-service/internal/config"
	"news-service/internal/repository/postgres"
	"news-service/internal/service"
	"news-service/internal/transport/grpc"
	"news-service/pkg/database"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadDefault()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация кеша
	cacheInstance := cache.New(cfg.Cache.TTL)
	defer cacheInstance.Stop()

	// Инициализация репозитория
	newsRepo := postgres.NewNewsRepository(db)

	// Инициализация сервиса
	newsService := service.NewNewsService(newsRepo, cacheInstance)

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer(newsService)

	// Канал для обработки сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
		log.Printf("Starting gRPC server on port %d", cfg.Server.GRPCPort)

		if err := grpcServer.Start(address); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Shutting down server...")

	// Graceful shutdown
	grpcServer.Stop()
	log.Println("Server stopped")
}
