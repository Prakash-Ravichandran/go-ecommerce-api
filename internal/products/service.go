package products

import (
	"context"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	ListProductsByID(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) ListProductsByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.ListProductsByID(ctx, id)
}
