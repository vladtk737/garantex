FROM golang:1.23.0-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY --from=builder /app/db/migrations /app/db/migrations
COPY .env /app/.env
CMD ["./app"]
