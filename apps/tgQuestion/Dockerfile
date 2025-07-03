FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY apps/tgQuestion/go.mod apps/tgQuestion/go.sum ./apps/tgQuestion/

RUN go mod download

COPY . .

WORKDIR /app/apps/tgQuestion
RUN go build -o /app/tg-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tg-service .

EXPOSE 50051

CMD ["./tg-service"]