# 📰 Микросервис новостей

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=for-the-badge&logo=google&logoColor=white)](https://grpc.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

> **Высокопроизводительный микросервис для управления новостями с gRPC API, PostgreSQL и in-memory кешированием**

---

## 🚀 Особенности

- ✅ **gRPC API** с полным CRUD функционалом
- ⚡ **In-memory кеш** с настраиваемым TTL
- 🗄️ **PostgreSQL** с миграциями и индексами
- 🐳 **Docker Compose** для быстрого развертывания
- 🧪 **Юнит-тесты** и graceful shutdown
- 📝 **Proto-спецификация** и автогенерация кода
- 🔧 **Makefile** с удобными командами

---

## 📋 Содержание

- [Быстрый старт](#-быстрый-старт)
- [API](#-api)
- [Архитектура](#-архитектура)
- [Конфигурация](#-конфигурация)
- [Команды](#-команды)
- [Разработка](#-разработка)
- [Тестирование](#-тестирование)
- [Деплой](#-деплой)

---

## ⚡ Быстрый старт

### Предварительные требования

- **Go** 1.21+
- **Docker** & **Docker Compose**
- **Make** (опционально)

### 🔧 Установка

```bash
# 1. Клонируйте репозиторий
git clone <repository-url>
cd news-service

# 2. Запустите PostgreSQL
make docker-up

# 3. Примените миграции
make migrate-up

# 4. Заполните тестовыми данными (опционально)
make seed

# 5. Запустите сервер
make run
```

**🎉 Готово!** Сервер доступен на `localhost:8080`

---

## 📚 API

### Модель новости

```protobuf
message News {
  string slug = 1;        // Уникальный идентификатор
  string title = 2;       // Заголовок (до 500 символов)
  string content = 3;     // Содержимое
  int64 created_at = 4;   // Время создания (Unix timestamp)
  int64 updated_at = 5;   // Время обновления (Unix timestamp)
}
```

### gRPC методы

| Метод | Описание | Кеширование |
|-------|----------|-------------|
| `CreateNews` | Создание новости | ➕ Добавляет в кеш |
| `GetNews` | Получение по slug | 🔍 Читает из кеша |
| `GetNewsList` | Список с пагинацией | 🔍 Кеширует списки |
| `UpdateNews` | Обновление по slug | 🔄 Инвалидирует кеш |
| `DeleteNews` | Удаление по slug | ❌ Удаляет из кеша |

### Примеры использования

```bash
# Создание новости
grpcurl -plaintext -d '{
  "slug": "my-news",
  "title": "Заголовок новости",
  "content": "Содержимое новости"
}' localhost:8080 news.NewsService/CreateNews

# Получение новости
grpcurl -plaintext -d '{"slug": "my-news"}' \
  localhost:8080 news.NewsService/GetNews

# Список новостей с пагинацией
grpcurl -plaintext -d '{"page": 1, "limit": 10}' \
  localhost:8080 news.NewsService/GetNewsList
```

---

## 🏗️ Архитектура

### Компоненты системы

```
┌─────────────────────────────────────────┐
│           gRPC API Layer                │
│      (Transport/Presentation)           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Service Layer                │
│         (Business Logic)                │
└─────────────────┬───────────────────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
      ▼           ▼           ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ In-Memory│ │PostgreSQL│ │   Cache  │
│   Cache  │ │Repository│ │ Manager  │
│          │ │          │ │          │
└──────────┘ └──────────┘ └──────────┘
```

### Структура проекта

```
news-service/
├── 📁 cmd/                    # Точки входа
│   ├── server/               # Основной сервер
│   ├── migrate/              # Команда миграций
│   └── seed/                 # Заполнение данными
├── 📁 internal/              # Приватный код
│   ├── cache/               # In-memory кеш
│   ├── config/              # Конфигурация
│   ├── domain/              # Доменные модели
│   ├── repository/          # Слой данных
│   ├── service/             # Бизнес-логика
│   └── transport/grpc/      # gRPC транспорт
├── 📁 proto/                 # Protocol Buffers
├── 📁 migrations/            # SQL миграции
├── 📁 pkg/                   # Публичные утилиты
├── 🐳 docker-compose.yml     # PostgreSQL для разработки
├── ⚙️ config.yml             # Конфигурация
└── 📝 Makefile              # Команды сборки
```

---

## ⚙️ Конфигурация

### Файл `config.yml`

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: news_db
  sslmode: disable

server:
  grpc_port: 8080

cache:
  ttl: 5m  # Время жизни кеша
```

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `password` |
| `DB_NAME` | Имя базы данных | `news_db` |
| `GRPC_PORT` | Порт gRPC сервера | `8080` |
| `CACHE_TTL` | TTL кеша | `5m` |

---

## 🛠️ Команды

### Основные команды

```bash
# Сборка и запуск
make build              # Сборка проекта
make run                # Запуск сервера
make test               # Запуск тестов

# Управление PostgreSQL
make docker-up          # Запуск PostgreSQL
make docker-down        # Остановка PostgreSQL

# Миграции
make migrate-up         # Применить миграции
make migrate-down       # Откатить миграции

# Данные
make seed               # Заполнить тестовыми данными

# Разработка
make proto              # Генерация protobuf
make clean-build        # Полная пересборка
make dev-setup          # Настройка окружения
```

---

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
make test

# Конкретный пакет
go test ./internal/cache/

# С покрытием
go test -cover ./...

# Verbose режим
go test -v ./internal/cache/
```

### Структура тестов

- ✅ **Юнит-тесты** для кеша (`internal/cache/cache_test.go`)
- ✅ **Thread-safety тесты** для concurrent доступа
- ✅ **TTL тесты** для автоматической очистки

---

## 🔧 Разработка

### Добавление новых методов API

1. **Обновите proto-файл** (`proto/news/news.proto`)
2. **Сгенерируйте код**: `make proto`
3. **Реализуйте в сервисе** (`internal/service/news.go`)
4. **Добавьте в transport** (`internal/transport/grpc/server.go`)

### Создание миграций

```bash
# Создайте файлы в migrations/
# XXX_migration_name.up.sql - применение
# XXX_migration_name.down.sql - откат

# Примените миграцию
make migrate-up
```

### Установка protoc

**Windows (Chocolatey):**
```bash
choco install protoc
```

**Linux (Ubuntu):**
```bash
sudo apt-get install protobuf-compiler
```

**macOS (Homebrew):**
```bash
brew install protobuf
```

---

## 📊 Производительность

### In-Memory Cache

- **Thread-safe** операции с RWMutex
- **Автоматическая очистка** просроченных записей
- **TTL** настраивается в конфигурации
- **Инвалидация** при изменениях данных

### База данных

- **Индексы** для оптимизации запросов
- **Автоматические триггеры** для `updated_at`
- **Connection pooling** в PostgreSQL драйвере

---

