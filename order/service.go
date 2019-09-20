package order

import (
	"context"
	"github.com/rs/xid"
	"time"
)

type Service interface {
	CreateOrder(ctx context.Context, userID string, products []OrderedBook) (*Order, error)
	GetOrdersForUser(ctx context.Context, userID string) ([]Order, error)
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s orderService) CreateOrder(ctx context.Context, userID string, books []OrderedBook) (*Order, error) {
	o := &Order{
		ID:        xid.New().String(),
		CreatedAt: time.Now().UTC(),
		UserID:    userID,
		Books:     books,
	}
	o.TotalPrice = 0.0
	for _, p := range books {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.CreateOrder(ctx, *o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s orderService) GetOrdersForUser(ctx context.Context, userID string) ([]Order, error) {
	return s.repository.GetOrdersForUser(ctx, userID)
}
