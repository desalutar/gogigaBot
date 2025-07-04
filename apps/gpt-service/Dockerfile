FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY pkg/go.mod ./pkg/
COPY apps/gpt-service/go.mod ./apps/gpt-service/

RUN go mod download

COPY . .

WORKDIR /app/apps/gpt-service
RUN go build -o /app/gpt-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/gpt-service .

EXPOSE 50052

CMD ["./gpt-service"]