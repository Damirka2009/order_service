package main

import (
	"log"
	pb "master/pkg/api/test"
	"net"

	"google.golang.org/grpc"

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

	log.Printf("gRPC server logging on: ", cfg.NETWORK+":"+cfg.GRPC_PORT)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingInterceptor),
	)

	svc := service.NewService()
	h := handler.NewOrderHandler(svc)

	pb.RegisterOrderServiceServer(grpcServer, h)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
