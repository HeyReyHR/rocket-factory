# Rocket Factory

Курсовой проект с курса [«Микросервисы на Go»](https://olezhek28.courses).
Микросервисная платформа управления заказами ракетных деталей: 7 сервисов, gRPC + REST, Kafka, PostgreSQL, MongoDB, полный стек observability.

---

## Стек технологий

| Категория | Технологии |
|-----------|------------|
| **Язык** | Go 1.24 |
| **Коммуникация** | gRPC (Protobuf), REST (OpenAPI/ogen), Kafka (event-driven) |
| **Базы данных** | PostgreSQL 17, MongoDB 7, Redis 7 |
| **Observability** | Prometheus, Grafana, Jaeger, OpenTelemetry Collector, Elasticsearch + Kibana |
| **Инфраструктура** | Docker, Docker Compose, multi-stage builds |
| **Кодогенерация** | Buf (proto), ogen (OpenAPI), mockery (моки) |
| **CI/CD** | GitHub Actions (lint, unit/integration tests, coverage) |
| **Автоматизация** | Taskfile |
| **Качество кода** | golangci-lint, gofumpt, gci |

---

## Архитектура

```
                        ┌──────────────────────────┐
                        │     Order Service         │
                        │   REST :8080 / gRPC :50050│
                        └─────┬──────────┬──────────┘
                     gRPC     │          │    gRPC
                  ┌───────────┘          └───────────┐
                  ▼                                   ▼
       ┌──────────────────┐                ┌──────────────────┐
       │ Inventory Service │                │  Payment Service  │
       │   gRPC :50051     │                │   gRPC :50052     │
       │   (MongoDB)       │                │   (stateless)     │
       └──────────────────┘                └────────┬─────────┘
                                                     │ Kafka: OrderPaid
                                        ┌────────────┼────────────┐
                                        ▼            ▼            ▼
                                ┌──────────┐  ┌──────────┐  ┌──────────┐
                                │ Assembly  │  │   IAM    │  │Notifica- │
                                │ Service   │  │ Service  │  │  tion    │
                                │(Postgres) │  │(PG+Redis)│  │(consumer)│
                                └─────┬─────┘  └──────────┘  └──────────┘
                                      │
                                Kafka: ShipAssembled
```

**Потоки данных:**
- **Синхронные** — Order вызывает Inventory и Payment по gRPC
- **Асинхронные** — события `OrderPaid` и `ShipAssembled` через Kafka
- **Observability** — все сервисы экспортируют метрики, трейсы и логи через OpenTelemetry

---

## Структура проекта

```
rocket-factory/
├── order/              # Сервис заказов (REST + gRPC, PostgreSQL)
├── inventory/          # Каталог деталей (gRPC, MongoDB)
├── payment/            # Обработка платежей (gRPC, stateless)
├── assembly/           # Сборка ракет (Kafka consumer, PostgreSQL)
├── iam/                # Аутентификация (PostgreSQL + Redis)
├── notification/       # Уведомления (Kafka consumer)
├── platform/           # Общие утилиты (logger, closer, kafka, grpc, tracing, metrics)
├── shared/             # Proto-файлы, OpenAPI-спецификации, сгенерированный код
│   ├── proto/          #   .proto определения сервисов и событий
│   ├── api/            #   OpenAPI спецификации
│   └── pkg/            #   Сгенерированный Go-код (proto + ogen)
├── deploy/
│   ├── compose/        # Docker Compose файлы для каждого сервиса
│   │   ├── core/       #   Kafka, Prometheus, Grafana, Jaeger, ELK, OTel
│   │   ├── order/
│   │   ├── inventory/
│   │   ├── assembly/
│   │   ├── iam/
│   │   ├── payment/
│   │   └── notification/
│   └── env/            # Шаблоны переменных окружения
├── .github/workflows/  # CI/CD пайплайны
├── Taskfile.yml        # Команды автоматизации
├── go.work             # Go workspace (7 модулей)
└── .golangci.yml       # Конфигурация линтера
```

Каждый сервис следует единой структуре:
```
<service>/
├── cmd/<service>_service/main.go   # Точка входа
├── internal/
│   ├── app/            # Инициализация и wiring зависимостей
│   ├── config/         # Загрузка конфигурации
│   ├── service/        # Бизнес-логика
│   ├── repository/     # Слой данных
│   ├── api/            # gRPC/HTTP хендлеры
│   ├── client/         # Клиенты к другим сервисам
│   ├── converter/      # DTO ↔ Model конвертеры
│   └── metrics/        # Prometheus метрики
├── migrations/         # SQL-миграции (Goose)
├── tests/              # Интеграционные тесты
└── go.mod
```

---

## Key Paths

| Путь | Описание |
|------|----------|
| `shared/proto/` | Protobuf-определения (inventory, events) |
| `shared/api/order/v1/order.openapi.yaml` | OpenAPI-спецификация Order Service |
| `platform/pkg/` | Переиспользуемые пакеты (kafka, grpc, logger, tracing, closer) |
| `deploy/compose/core/` | Инфраструктура (Kafka, Prometheus, Grafana, Jaeger, ELK) |
| `order/migrations/` | SQL-миграции (Goose) |
| `.golangci.yml` | Конфигурация линтера |
| `deploy/env/.env.template` | Шаблон переменных окружения |

---

## Локальный запуск

### Требования

- Docker & Docker Compose
- Go 1.24+
- [Task](https://taskfile.dev/) (task runner)
- Node.js (для Redocly CLI)

### Запуск

```bash
# 1. Создать Docker-сеть (один раз)
docker network create microservices-net

# 2. Сгенерировать .env файлы из шаблонов
task env:generate

# 3. Поднять инфраструктуру (Kafka, Prometheus, Grafana, Jaeger, ELK)
task up-core

# 4. Поднять сервисы
task up-inventory
task up-order
task up-assembly
task up-iam

# Или всё сразу:
task up-all
```

### UI-панели

| Сервис | URL |
|--------|-----|
| Grafana | http://localhost:3000 (admin/admin) |
| Prometheus | http://localhost:9090 |
| Jaeger | http://localhost:16686 |
| Kibana | http://localhost:5601 |
| Kafka UI | http://localhost:8090 |
| Order API | http://localhost:8080 |

### Остановка

```bash
task down-all
```

---

## Task-команды

### Кодогенерация
| Команда | Описание |
|---------|----------|
| `task gen` | Генерация proto + OpenAPI кода |
| `task proto:gen` | Генерация Go-кода из .proto файлов |
| `task ogen:gen` | Генерация Go-кода из OpenAPI спецификаций |
| `task mockery:gen` | Генерация моков интерфейсов |

### Качество кода
| Команда | Описание |
|---------|----------|
| `task format` | Форматирование (gofumpt + gci) |
| `task lint` | Запуск golangci-lint |
| `task test` | Юнит-тесты |
| `task test-integration` | Интеграционные тесты (с testcontainers) |
| `task test-coverage` | Покрытие бизнес-логики |
| `task coverage:html` | HTML-отчёт покрытия |

### Инфраструктура
| Команда | Описание |
|---------|----------|
| `task up-core` | Kafka, Prometheus, Grafana, Jaeger, ELK |
| `task up-inventory` | Inventory + MongoDB |
| `task up-order` | Order + PostgreSQL |
| `task up-assembly` | Assembly + PostgreSQL |
| `task up-iam` | IAM + PostgreSQL + Redis |
| `task up-all` / `task down-all` | Все сервисы |
| `task env:generate` | Генерация .env из шаблонов |

### API-тестирование
| Команда | Описание |
|---------|----------|
| `task test-api` | E2E тест: создание заказа → оплата → сборка → отмена |

### Зависимости
| Команда | Описание |
|---------|----------|
| `task deps:update` | `go mod tidy` для всех модулей |
