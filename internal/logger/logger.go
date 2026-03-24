package logger

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		log.Printf("method=%s duration=%s status=ERROR error=%v", info.FullMethod, duration, err)
	} else {
		log.Printf("method=%s duration=%s status=OK", info.FullMethod, duration)
	}

	return resp, nil
}
