FROM golang:1.25-alpine AS build
WORKDIR /src

COPY Webby.AdminService/go.mod Webby.AdminService/go.sum ./
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
COPY --from=build /src/docs ./docs

EXPOSE 8085

ENTRYPOINT ["./app"]