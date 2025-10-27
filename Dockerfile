# syntax=docker/dockerfile:1.4

FROM --platform=$TARGETPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Отладка архитектуры
RUN echo "Building for: TARGETOS=$TARGETOS, TARGETARCH=$TARGETARCH, BUILDPLATFORM=$BUILDPLATFORM" && uname -m

# Сборка бинарей под целевую архитектуру
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/datagen ./db_datagen
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/venture-platform ./cmd

FROM alpine:latest

WORKDIR /root/
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/datagen /datagen
COPY --from=builder /app/venture-platform /venture-platform

RUN chmod +x /datagen /venture-platform

EXPOSE 8080
ENTRYPOINT ["/datagen"]
