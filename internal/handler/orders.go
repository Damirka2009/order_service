package handler

import (
	"context"
	"master/internal/service"
	pb "master/pkg/api/test"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	service *service.Service
}

func NewOrderHandler(s *service.Service) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id := h.service.Create(req.Item, req.Quantity)

	return &pb.CreateOrderResponse{
		Id: id,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.service.Get(req.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return &pb.GetOrderResponse{Order: order}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	if req.Item == "" || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}
	order, err := h.service.Update(req.Id, req.Item, req.Quantity)

	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return &pb.UpdateOrderResponse{
		Order: order,
	}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	_, err := h.service.Delete(req.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return &pb.DeleteOrderResponse{
		Success: true,
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders := h.service.List()

	return &pb.ListOrdersResponse{
		Orders: orders,
	}, nil
}
