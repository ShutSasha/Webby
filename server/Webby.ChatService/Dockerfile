FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.ChatService/go.mod Webby.ChatService/go.sum ./
RUN go mod download
COPY Webby.ChatService/ .
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.ChatService/config/config.yaml config/config.yaml
EXPOSE 8081 50051
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
