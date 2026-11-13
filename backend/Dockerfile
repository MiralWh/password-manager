FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.* ./
RUN go mod download || true
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata wget && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/server ./server
COPY config.prod.yaml ./config.yaml
USER app
EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --start-period=10s --retries=10 \
  CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["./server", "-config", "/app/config.yaml"]
