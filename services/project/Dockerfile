FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /project_service ./cmd

FROM alpine:3.18
COPY --from=build /project_service /project_service
EXPOSE 8104
ENTRYPOINT ["/project_service"]
