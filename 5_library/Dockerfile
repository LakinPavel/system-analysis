FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o library-service ./cmd/library
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/library-service .
COPY --from=builder /app/db/migrations ./db/migrations
EXPOSE 8080 8081
CMD ["./library-service"]