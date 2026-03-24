package main

import (
	"context"
	"log"
	pb "master/pkg/api/test"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"master/internal/config"
	"master/internal/handler"
	"master/internal/logger"
	"master/internal/service"
)

func main() {
	cfg := config.Load()
	lis, err := net.Listen(cfg.NETWORK, ":"+cfg.GRPC_PORT)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingInterceptor),
	)

	svc := service.NewService()
	h := handler.NewOrderHandler(svc)

	pb.RegisterOrderServiceServer(grpcServer, h)
	go func() {
		log.Printf("gRPC server logging on: %s", cfg.NETWORK+":"+cfg.GRPC_PORT)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	err = pb.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		cfg.NETWORK+":"+cfg.GRPC_PORT,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP_PORT,
		Handler: mux,
	}

	go func() {
		log.Println("HTTP gateway listening on:" + cfg.HTTP_PORT)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down servers....")

	ctxShutdown, cancelShotdown := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelShotdown()

	if err := httpServer.Shutdown(ctxShutdown); err != nil {
		log.Fatalln("HTTP shutdown error: %v", err)
	}

	grpcServer.GracefulStop()
	log.Println("Servers stopped gracefully")
}
