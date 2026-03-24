package service

import (
	"errors"
	"fmt"
	pb "master/pkg/api/test"
	"sync"
	"sync/atomic"
)

var (
	ErrNotFound = errors.New("order not found")
)

type Service struct {
	mu        sync.RWMutex
	orders    map[string]*pb.Order
	idCounter uint64
}

func NewService() *Service {
	return &Service{
		orders: make(map[string]*pb.Order),
	}
}

func (s *Service) Create(item string, quantity int32) string {
	id := atomic.AddUint64(&s.idCounter, 1)
	orderID := fmt.Sprintf("order-%d", id)

	order := &pb.Order{
		Id:       orderID,
		Item:     item,
		Quantity: quantity,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[orderID] = order

	return orderID
}

func (s *Service) Get(id string) (*pb.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	if !ok {
		return nil, ErrNotFound
	}

	return order, nil
}

func (s *Service) Update(id string, item string, quantity int32) (*pb.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return nil, ErrNotFound
	}

	order.Item = item
	order.Quantity = quantity
	return order, nil
}

func (s *Service) Delete(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[id]
	if !ok {
		return false, ErrNotFound
	}
	delete(s.orders, id)
	return true, nil
}

func (s *Service) List() []*pb.Order {
	list := make([]*pb.Order, 0, len(s.orders))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, order := range s.orders {
		list = append(list, order)
	}

	return list
}
