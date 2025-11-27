FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /comment_service ./cmd

FROM alpine:3.18
COPY --from=build /comment_service /comment_service
EXPOSE 8105
ENTRYPOINT ["/comment_service"]
