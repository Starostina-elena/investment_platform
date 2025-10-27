FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию.
WORKDIR /app

# Копируем go.mod и go.sum, чтобы кэшировать зависимости.
# Это, блядь, важно для скорости пересборки!
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код.
COPY . .

# Собираем наше приложение. CGO_ENABLED=0 создает статический бинарник без C-зависимостей.
# -o /datagen - имя выходного файла.
RUN CGO_ENABLED=0 GOOS=linux go build -o /datagen .


# --- ЭТАП 2: ФИНАЛЬНЫЙ ОБРАЗ (The Production Commune) ---
# Начинаем с чистого, минималистичного образа. Никаких, нахуй, компиляторов!
FROM alpine:latest

# Копируем ТОЛЬКО скомпилированный бинарник из этапа сборщика.
COPY --from=builder /datagen /datagen

# Указываем, что при запуске контейнера нужно выполнить наш бинарник.
ENTRYPOINT ["/datagen"]