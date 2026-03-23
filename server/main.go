package main

import (
	"context"
	"fmt"
	"log"
	pb "master/pkg/api/test"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	orders    map[string]*pb.Order
	mu        sync.RWMutex
	idCounter uint64
}

func newServer() *server {
	return &server{
		orders: make(map[string]*pb.Order),
	}
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id := atomic.AddUint64(&s.idCounter, 1)
	orderID := fmt.Sprintf("order-%d", id)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[orderID] = &pb.Order{
		Id:       orderID,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	return &pb.CreateOrderResponse{
		Id: orderID,
	}, nil
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return &pb.GetOrderResponse{
		Order: order,
	}, nil
}

func (s *server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	s.orders[req.Id] = &pb.Order{
		Item:     req.Item,
		Quantity: req.Quantity,
		Id:       req.Id,
	}
	return &pb.UpdateOrderResponse{
		Order: s.orders[req.Id],
	}, nil
}

func (s *server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	delete(s.orders, req.Id)
	return &pb.DeleteOrderResponse{
		Success: true,
	}, nil
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*pb.Order, 0, len(s.orders))
	for _, order := range s.orders {
		list = append(list, order)
	}
	return &pb.ListOrdersResponse{
		Orders: list,
	}, nil
}

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

func main() {
	// .env
	godotenv.Load()
	grpc_port := os.Getenv("GRPC_PORT")
	network := os.Getenv("Network")

	lis, err := net.Listen(network, ":"+grpc_port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
	)

	srv := newServer()
	pb.RegisterOrderServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
