package products

import (
	"context"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	ListProductsByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProducts(ctx context.Context, product repo.CreateProductParams) (repo.Product, error)
	UpdateProductPrice(ctx context.Context, product repo.UpdateProductPriceParams) (repo.Product, error)
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

func (s *svc) CreateProducts(ctx context.Context, product repo.CreateProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, product)
}

func (s *svc) UpdateProductPrice(ctx context.Context, product repo.UpdateProductPriceParams) (repo.Product, error) {
	return s.repo.UpdateProductPrice(ctx, product)
}
