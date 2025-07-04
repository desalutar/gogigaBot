FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY apps/tg-service/go.mod apps/tg-service/go.sum ./apps/tg-service/

RUN go mod download

COPY . .

WORKDIR /app/apps/tg-service
RUN go build -o /app/tg-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tg-service .

EXPOSE 50051

CMD ["./tg-service"]