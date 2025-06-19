.PHONY: build run test proto proto-docker migrate-up migrate-down seed docker-up docker-down

# Генерация protobuf файлов (локально)
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/news/news.proto

# Генерация protobuf файлов через Docker (если нет локальных плагинов)
proto-docker:
	docker run --rm -v $(PWD):/workspace -w /workspace \
		namely/protoc-all:1.51_1 \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/news/news.proto

# Установка protoc плагинов
install-proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Сборка проекта
build:
	go build -o bin/server cmd/server/main.go
	go build -o bin/migrate cmd/migrate/main.go
	go build -o bin/seed cmd/seed/main.go

# Запуск сервера
run:
	go run cmd/server/main.go

# Тесты
test:
	go test ./...

# Поднятие PostgreSQL
docker-up:
	docker-compose up -d

# Остановка PostgreSQL
docker-down:
	docker-compose down

# Применение миграций
migrate-up:
	go run cmd/migrate/main.go up

# Откат миграций
migrate-down:
	go run cmd/migrate/main.go down

# Заполнение тестовыми данными
seed:
	go run cmd/seed/main.go

# Полная пересборка
clean-build: docker-down docker-up install-proto-deps proto build

# Быстрый старт для разработки
dev-setup: install-proto-deps docker-up migrate-up seed
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the server"