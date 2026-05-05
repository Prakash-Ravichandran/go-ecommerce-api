package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type OrderService interface {
	GetOrders(ctx context.Context) ([]repo.Order, error)
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product has not enough stock")
)

func NewService(repo *repo.Queries, db *pgx.Conn) OrderService {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) GetOrders(ctx context.Context) ([]repo.Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	// validate payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}
	// create an order
	// look for the product if exits
	// create order item

	//create a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	// if there is an err then rollback the changes
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// look for the product if exits
	for _, item := range tempOrder.Items {
		product, err := qtx.ListProductsByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, nil
		}
		// challenge: update the product stock quantity
	}

	tx.Commit(ctx) // save order after creating it, if not changes won't be saved to DB.

	return order, nil
}
