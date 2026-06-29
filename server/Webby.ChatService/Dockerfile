FROM golang:1.25-alpine AS build
WORKDIR /src

COPY Webby.ChatService/go.mod Webby.ChatService/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=chat-go-mod \
    go mod download

COPY Webby.ChatService/ .
RUN --mount=type=cache,target=/go/pkg/mod,id=chat-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=chat-go-build \
    CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /app .
COPY --from=build /src/docs ./docs

EXPOSE 8081 50051

ENTRYPOINT ["./app"]