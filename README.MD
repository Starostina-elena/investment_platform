# Venture Platform - Backend

Проект представляет собой микросервисную платформу для управления инвестиционными проектами. Организации и физические лица, заинтересованные в сборе средств, могут создавать сборы, а пользователи — вносить пожертвования в выбранные проекты. В зависимости от типа сбора (благотворительность, фиксированный процент от суммы, процентная ставка [долг растёт со временем]) после использования средств организацией система выполняет автоматический расчёт возврата.


В проекте используется тестовый API ЮКасса для имитации ввода и вывода средств пользователями. Предусмотрена загрузка документов организациями (хранение в MinIO), кэширование части данных в Redis для ускорения работы приложения и имитация отправки уведомлений по email с помощью Mailhog.


При разработке системы сделан упор на безопасность и производительность. Реализована система модерации: администратор может ограничивать действия пользователя, организации и проекта [ сбора ]. Реализован контроль уровней доступа пользователей, в частности, организация может добавлять сотрудников с различными правами. Тестирование с помощью Apache JMeter показало пропускную способность ~1900 RPS на эндпоинте выдачи проектов главной страницы.

## Документация
Подробное описание реализованных эндпоинтов доступно по [ссылке](https://docs.google.com/document/d/1vLhRPkJyZQ2e4wBFpFiV9_42txXgDAXzqB6JgmLHKjc/edit?usp=sharing).

## Архитектура

Бэкенд состоит из набора независимых сервисов, разворачиваемых через `docker-compose` и объединённых в одну сеть `app-network`:

- User (`services/user`, порт 8101): пользователи, профили, аватары (MinIO).
- Organisation (`services/organisation`, порт 8102): организации, документы (MinIO, `orgdocs`).
- Transactions (`services/transactions`, порт 8103): транзакции, интеграция с уведомлениями.
- Project (`services/project`, порт 8104): проекты, медиа (MinIO `projects`), зависимости от Organisation и Transactions.
- Comment (`services/comment`, порт 8105): комментарии к проектам.
- Payment (`services/payment`, порт 8106): платежи и выводы, интеграция с YooKassa, взаимодействует с Transactions.
- Notification (`services/notification`, порт 8083): отправка email (SMTP), используется другими сервисами.
- Daemon (`services/daemon`): фоновые задачи, использует БД и Notification.
- Gateway (`nginx`, порт 80): единая точка входа, проксирует запросы на микросервисы.
- База данных: PostgreSQL 15 (порт 5432).
- Объектное хранилище: MinIO (порты 9000/9001).
- Кеш/очереди: Redis (порт 6379).
- Mailhog (порты 1025 SMTP / 8025 Web UI) для разработки.

Также присутствует контейнер `app` (порт 8080) со сборкой двоичных файлов:
- `venture-platform` — приложение, которое подключается к БД и применяет миграции (см. `cmd/main.go`).
- `datagen` — генератор данных (`db_datagen`), может применяться для первичного наполнения.

## Зависимости

- Docker или Podman (для `compose`)
- Go 1.25 
- Postgres 15, Redis, MinIO, Mailhog 

## Конфигурация окружения

`docker-compose.yml` использует переменные окружения (подставляются из `.env` в корне проекта)

```dotenv
# YooKassa
YOOKASSA_SHOP_ID=your_shop_id
YOOKASSA_SECRET_KEY=your_secret_key
YOOKASSA_AGENT_ID=your_agent_id
YOOKASSA_PAYOUT_API_KEY=your_payout_api_key

```

Сервисные переменные для подключения к БД и MinIO уже заданы внутри `docker-compose.yml`:

- БД: `DB_HOST=db`, `DB_PORT=5432`, `DB_USER=comrade`, `DB_PASSWORD=secret_password`, `DB_NAME=venture-platform-db`.
- MinIO: `MINIO_ENDPOINT=minio:9000`, `MINIO_ACCESS_KEY=minioadmin`, `MINIO_SECRET_KEY=minioadmin123`, `MINIO_USE_SSL=false`.

## Миграции базы данных

Миграции хранятся в каталоге `db/migrations` и применяются автоматически приложением `venture-platform` при старте:

- Внутри `cmd/main.go` приложение подключается к БД, затем вызывает `db.Migrate()`
- Миграции встроены через `embed` (см. `db/migrate.go`) и применяются с помощью `golang-migrate`

## Быстрый старт (Docker Compose)

1. Установите Docker Desktop
2. Создайте файл `.env` в корне, заполните по образцу с .env.example
3. Запустите все сервисы:

```bash
make up
```

Команда `make` использует доступный рантайм (`docker` или `podman`) и выполнит `compose up --build -d`.

Остановить контейнеры:

```bash
make down
```

### Порты и доступ

- Gateway (Nginx): http://localhost:80
- Frontend (Next.js): http://localhost:3000
- User: http://localhost:12550
- Organisation: http://localhost:12551
- Transactions: http://localhost:12552
- Project: http://localhost:12553
- Comment: http://localhost:12554
- Notification: http://localhost:12555
- Payment: http://localhost:12556
- App (venture-platform): http://localhost:8080
- Postgres: localhost:5432
- MinIO: консоль http://localhost:9001, API http://localhost:9000
- Redis: localhost:6379
- Mailhog: SMTP localhost:1025, UI http://localhost:8025
