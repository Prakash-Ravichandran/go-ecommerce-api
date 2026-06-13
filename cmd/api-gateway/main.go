package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayHandler struct {
	productClient pbProduct.ProductServiceClient
}

const gatewayPortNum string = ":8080"

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Establish a gRPC connection to your existing Product Service (:50051)
	productServiceAddr := ":50051"
	conn, err := grpc.Dial(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Gateway failed to connect to Product Service", "error", err)
		os.Exit(1)
	}

	defer conn.Close()

	client := pbProduct.NewProductServiceClient(conn)
	handler := &GatewayHandler{productClient: client}

	// 2. setup standard HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/products", handler.handleProducts)

	log.Println("starting our api gateway server")

	log.Println("Started on port", gatewayPortNum)
	fmt.Println("To close connection CTRL+C :-)")

	// spinning up the server
	gatewaySrvErr := http.ListenAndServe(gatewayPortNum, mux)
	if gatewaySrvErr != nil {
		slog.Error("API Gateway crashed", "error", gatewaySrvErr)
		log.Fatal(gatewaySrvErr)
	}
}

func (h *GatewayHandler) handleProducts(w http.ResponseWriter, r *http.Request) {
	slog.Info("hanlder")
}
