FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment_service ./cmd

FROM alpine:3.18
COPY --from=build /payment_service /payment_service
EXPOSE 8106
ENTRYPOINT ["/payment_service"]