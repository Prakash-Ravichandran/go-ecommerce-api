package main

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/health"
	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	pb.UnimplementedHealthServiceServer
}

type grpcServerForProducts struct {
	pbProduct.UnimplementedProductServiceServer
}

func (s *grpcServer) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	slog.Info("gRPC health invoked")

	return &pb.HealthResponse{
		Status: "all good",
	}, nil
}

func (s *grpcServerForProducts) GetProducts(ctx context.Context, req *pbProduct.GetProductsRequest) (*pbProduct.GetProductsResponse, error) {
	slog.Info("gRPC GetProducts invoked")

	return &pbProduct.GetProductsResponse{
		Id:           12,
		Name:         "Macbook4",
		PriceInCents: 55,
		Quantity:     45,
		CreatedAt:    timestamppb.New(time.Now()),
	}, nil
}
