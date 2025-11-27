FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /organisation_service ./cmd

FROM alpine:3.18
COPY --from=build /organisation_service /organisation_service
EXPOSE 8102
ENTRYPOINT ["/organisation_service"]
