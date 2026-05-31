FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.WsGateway/go.mod Webby.WsGateway/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=wsgateway-go-mod \
    go mod download
COPY Webby.WsGateway/ .
RUN --mount=type=cache,target=/go/pkg/mod,id=wsgateway-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=wsgateway-go-build \
    CGO_ENABLED=0 go build -o /app ./cmd/server
    
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.WsGateway/config/config.yaml config/config.yaml
COPY Webby.WsGateway/docs/ docs/
EXPOSE 8090
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
