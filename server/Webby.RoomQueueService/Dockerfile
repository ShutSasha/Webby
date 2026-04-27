FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.RoomQueueService/go.mod Webby.RoomQueueService/go.sum ./
RUN go mod download
COPY Webby.RoomQueueService/ .
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.RoomQueueService/config/config.yaml config/config.yaml
COPY Webby.RoomQueueService/migrations/ migrations/
EXPOSE 8083 50053
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
