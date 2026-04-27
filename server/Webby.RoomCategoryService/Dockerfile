FROM golang:1.25-alpine AS build
WORKDIR /src
COPY Webby.RoomCategoryService/go.mod Webby.RoomCategoryService/go.sum ./
RUN go mod download
COPY Webby.RoomCategoryService/ .
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app .
COPY Webby.RoomCategoryService/config/config.yaml config/config.yaml
EXPOSE 8082 50052
ENV CONFIG_PATH=config/config.yaml
ENTRYPOINT ["./app"]
