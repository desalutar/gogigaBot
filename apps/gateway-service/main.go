package main

import (
	"context"
	"gptBot/pkg/gen/gpt-service"
	"log"
	"net/http"

	// "runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }

    err := gpt.RegisterQAServiceHandlerFromEndpoint(ctx, mux, "gpt-service:50052", opts)
    if err != nil {
        log.Fatalf("Failed to register gRPC gateway: %v", err)
    }

    log.Println("Starting HTTP Gateway on :8080")
    err = http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatalf("gateway failed: %v", err)
    }
}