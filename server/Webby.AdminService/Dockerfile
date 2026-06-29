FROM golang:1.25-alpine AS build
WORKDIR /src

COPY Webby.AdminService/go.mod Webby.RoomCategoryService/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=admin-go-mod \
    go mod download

COPY Webby.AdminService/ .
RUN --mount=type=cache,target=/go/pkg/mod,id=admin-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=admin-go-build \
    CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.AdminService/config/config.yaml config/config.yaml
EXPOSE 8082 50052
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
