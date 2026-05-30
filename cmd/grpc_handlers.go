package main

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/health"
	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type grpcServer struct {
	pb.UnimplementedHealthServiceServer
	pbProduct.UnimplementedProductServiceServer
	db   *pgx.Conn
	repo repo.Querier
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
	slog.Info("Checking gRPC dependencies",
		"is_s_nil", s == nil,
		"is_db_nil", s.db == nil,
		"is_repo_nil", s.repo == nil,
	)

	if s.repo == nil {
		slog.Error("CRITICAL: s.repo is completely nil inside handler!")
		return nil, context.DeadlineExceeded // return an explicit error to avoid panic
	}

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

	return &pbProduct.UpdateProductsResponse{
		Product: &pbProduct.Product{
			Id:           12,
			Name:         "Macbook4",
			PriceInCents: 55,
			Quantity:     45,
			CreatedAt:    timestamppb.New(time.Now()),
		},
	}, nil
}
