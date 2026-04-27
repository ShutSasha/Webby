FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.RoomService/go.mod Webby.RoomService/go.sum ./
RUN go mod download
COPY Webby.RoomService/ .
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.RoomService/config/config.yaml config/config.yaml
COPY Webby.RoomService/migrations/ migrations/
EXPOSE 8080 50054
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
