package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	pb "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/health"
	pbOrders "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/orders"
	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type grpcServer struct {
	pb.UnimplementedHealthServiceServer
	pbProduct.UnimplementedProductServiceServer
	pbOrders.UnimplementedOrderServiceServer
	db   *pgx.Conn
	repo repo.Queries
}

func (s *grpcServer) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	slog.Info("gRPC health invoked")

	return &pb.HealthResponse{
		Status: "all good",
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pbProduct.GetProductsRequest) (*pbProduct.GetProductsResponse, error) {
	slog.Info("gRPC GetProducts invoked")

	// 1. DEBUG LOGS: Check what is actually nil before running the query
	// slog.Info("Checking gRPC dependencies",
	// 	"is_s_nil", s == nil,
	// 	"is_db_nil", s.db == nil,
	// 	"is_repo_nil", s.repo == nil,
	// )

	// if s.repo == nil {
	// 	slog.Error("CRITICAL: s.repo is completely nil inside handler!")
	// 	return nil, context.DeadlineExceeded // return an explicit error to avoid panic
	// }

	dbproducts, err := s.repo.ListProducts(ctx)

	if err != nil {
		slog.Info("Error Fetching Products")
	}

	var grpcProducts []*pbProduct.Product

	for _, p := range dbproducts {
		grpcProduct := &pbProduct.Product{
			Id:           p.ID,
			Name:         p.Name,
			PriceInCents: int32(p.PriceInCents),
			Quantity:     int32(p.Quantity),
			CreatedAt:    timestamppb.New(p.CreatedAt.Time),
		}
		grpcProducts = append(grpcProducts, grpcProduct)
	}

	return &pbProduct.GetProductsResponse{
		Products: grpcProducts,
	}, nil
}

func (s *grpcServer) ListProductsByID(ctx context.Context, req *pbProduct.ListProductsByIDRequest) (*pbProduct.ListProductsByIDResponse, error) {
	slog.Info("gRPC ListProductsByID invoked")

	tempListProductById := req.GetId()

	product, err := s.repo.ListProductsByID(ctx, tempListProductById)

	if err != nil {
		slog.Error("DataBase failed to get product by id")
	}

	slog.Info("Product Fetched successfully")

	return &pbProduct.ListProductsByIDResponse{
		Product: &pbProduct.Product{
			Id:           product.ID,
			Name:         product.Name,
			PriceInCents: product.PriceInCents,
			Quantity:     product.Quantity,
			CreatedAt:    timestamppb.New(product.CreatedAt.Time),
		},
	}, nil
}

func (s *grpcServer) CreateProducts(ctx context.Context, req *pbProduct.CreateProductsRequest) (*pbProduct.CreateProductsResponse, error) {

	slog.Info("create products invoked")

	now := time.Now()
	pgTime := pgtype.Timestamp{
		Time:  now,
		Valid: true,
	}

	tempProduct := repo.CreateProductParams{
		ID:           req.GetId(),
		Name:         req.GetName(),
		Quantity:     req.GetQuantity(),
		PriceInCents: req.GetPriceInCents(),
		CreatedAt:    pgtype.Timestamptz(pgTime),
	}

	createdProduct, err := s.repo.CreateProduct(ctx, tempProduct)

	if err != nil {
		slog.Error("Database failed to insert product", "error", err)
		return nil, err // Return the actual database failure to prevent crashing below
	}

	slog.Info("Product successfully created in DB", "id", createdProduct.ID)

	return &pbProduct.CreateProductsResponse{
		Product: &pbProduct.Product{
			Id:           createdProduct.ID,
			Name:         createdProduct.Name,
			PriceInCents: createdProduct.PriceInCents,
			Quantity:     createdProduct.Quantity,
			CreatedAt:    timestamppb.New(createdProduct.CreatedAt.Time),
		},
	}, nil
}

func (s *grpcServer) UpdateProducts(ctx context.Context, req *pbProduct.UpdateProductsRequest) (*pbProduct.UpdateProductsResponse, error) {
	slog.Info("update products invoked")

	tempUpdateProduct := repo.UpdateProductPriceParams{
		ID:           req.GetId(),
		PriceInCents: req.GetPriceInCents(),
	}

	updatedProduct, err := s.repo.UpdateProductPrice(ctx, tempUpdateProduct)

	if err != nil {
		slog.Error("Database failed to update product proce")
	}

	slog.Info("Product Price updated successfully")

	return &pbProduct.UpdateProductsResponse{
		Product: &pbProduct.Product{
			Id:           updatedProduct.ID,
			Name:         updatedProduct.Name,
			PriceInCents: updatedProduct.PriceInCents,
			Quantity:     updatedProduct.Quantity,
			CreatedAt:    timestamppb.New(time.Now()),
		},
	}, nil
}

// orders handlers
func (s *grpcServer) GetOrders(ctx context.Context, req *pbOrders.GetOrdersRequest) (*pbOrders.GetOrdersResponse, error) {
	slog.Info("gRPC GetOrders invoked")

	orders, err := s.repo.ListOrders(ctx)
	if err != nil {
		slog.Error("Error fetching orders")
	}

	var grpcOrders []*pbOrders.Order

	for _, o := range orders {
		grpcOrder := &pbOrders.Order{
			Id:         o.ID,
			CustomerId: o.CustomerID,
			CreatedAt:  timestamppb.New(o.CreatedAt.Time),
		}

		grpcOrders = append(grpcOrders, grpcOrder)
	}

	return &pbOrders.GetOrdersResponse{
		Orders: grpcOrders,
	}, nil
}

func (s *grpcServer) CreateOrders(ctx context.Context, req *pbOrders.CreateOrdersRequest) (*pbOrders.CreateOrdersResponse, error) {
	slog.Info("gRPC CreateOrders invoked")

	// validate payload from the incoming request
	if req.GetCustomerId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}
	if len(req.GetItems()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "at least one item is required")
	}
	// create an order
	// look for the product if exits
	// create order item

	//create a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	// if there is an err then rollback the changes
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// 1. Look for an existing order
	existingOrders, err := qtx.ListOrdersByCustomerID(ctx, req.GetCustomerId())

	// If the database tells us specifically "no rows found", that is GREAT news.
	// It means the customer is new, so we clear the error and proceed.
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		existingOrders = nil
	}

	// If there was a real database error (connection drop, syntax issue), halt here
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database lookup failed: %v", err)
	}

	// Now check if we actually found a record
	// This works whether existingOrders is a slice (len > 0) or a single struct entity checking an ID match
	if len(existingOrders) > 0 {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"customer %d already has an active order; duplicate orders are blocked",
			req.GetCustomerId(),
		)
	}

	//create an order
	order, err := qtx.CreateOrder(ctx, req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// look for the product if exits
	for _, item := range req.GetItems() {
		product, err := qtx.ListProductsByID(ctx, item.ProductId)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "product with ID %d not found", item.ProductId)
		}

		if product.Quantity < int32(item.Quantity) {
			return nil, status.Errorf(codes.FailedPrecondition, "product %s has insufficient stock", product.Name)
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductId,
			Quantity:   int32(item.Quantity),
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
		}
		// challenge: update the product stock quantity
	}

	tx.Commit(ctx) // save order after creating it, if not changes won't be saved to DB.

	return &pbOrders.CreateOrdersResponse{
		Order: &pbOrders.Order{
			Id:         order.ID,
			CustomerId: order.CustomerID,
			CreatedAt:  timestamppb.New(order.CreatedAt.Time),
		},
	}, nil
}
