FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY apps/gateway-service/go.mod apps/gateway-service/go.sum ./apps/gateway-service/

RUN go mod download

COPY . .

WORKDIR /app/apps/gateway-service
RUN go build -o /app/gateway-service ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/gateway-service .

EXPOSE 8080

CMD ["./gateway-service"]