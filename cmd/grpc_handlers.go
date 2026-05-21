package main

import (
	"context"
	"log/slog"

	pb "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/health"
)

type grpcServer struct {
	pb.UnimplementedHealthServiceServer
}

func (s *grpcServer) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	slog.Info("gRPC health invoked")

	return &pb.HealthResponse{
		Status: "all good",
	}, nil
}
