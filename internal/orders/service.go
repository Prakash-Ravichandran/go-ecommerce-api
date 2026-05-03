package orders

import (
	"context"
	"time"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService interface {
	GetOrder(ctx context.Context) []repo.Order
}

type svc struct{}

func NewService() OrderService {
	return &svc{}
}

func (s *svc) GetOrder(ctx context.Context) []repo.Order {
	dummyOrders := []repo.Order{
		{
			ID: 1, CustomerID: 12, CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}},
		{
			ID: 2, CustomerID: 49, CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}},
		{
			ID: 3, CustomerID: 79, CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}},
	}
	return dummyOrders
}
