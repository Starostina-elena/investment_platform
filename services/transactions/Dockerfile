FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /transactions_service ./cmd

FROM alpine:3.18
COPY --from=build /transactions_service /transactions_service
EXPOSE 8103
ENTRYPOINT ["/transactions_service"]
