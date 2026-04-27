FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.VotesService/go.mod Webby.VotesService/go.sum ./
RUN go mod download
COPY Webby.VotesService/ .
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.VotesService/config/config.yaml config/config.yaml
COPY Webby.VotesService/migrations/ migrations/
EXPOSE 8084 50055
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
