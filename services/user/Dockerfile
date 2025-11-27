FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /user_service ./cmd

FROM alpine:3.18
COPY --from=build /user_service /user_service
EXPOSE 8101
ENTRYPOINT ["/user_service"]
