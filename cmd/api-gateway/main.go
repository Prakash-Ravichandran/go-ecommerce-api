package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	pbOrders "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/orders"
	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GatewayHandler struct {
	productClient pbProduct.ProductServiceClient
	ordersClient  pbOrders.OrderServiceClient
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

	orderClient := pbOrders.NewOrderServiceClient(conn)
	ordersHandler := &GatewayHandler{ordersClient: orderClient}

	// 2. setup standard HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/products", handler.handleProducts)
	mux.HandleFunc("/orders", ordersHandler.handleOrders)

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
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		slog.Info("Gateway: GET /products intercepted. Forwarding to gRPC :50051")
		res, err := h.productClient.GetProducts(ctx, &pbProduct.GetProductsRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(res)

	case http.MethodPost:
		slog.Info("Gateway: POST /products intercepted. Forwarding to gRPC :50051")

		var protoReq pbProduct.CreateProductsRequest
		if err := json.NewDecoder(r.Body).Decode(&protoReq); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		res, err := h.productClient.CreateProducts(ctx, &protoReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h *GatewayHandler) handleOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		slog.Info("Gateway: GET /orders intercepted. Forwarding to gRPC :50052")
		res, err := h.ordersClient.GetOrders(ctx, &pbOrders.GetOrdersRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(res)

	case http.MethodPost:
		slog.Info("Gateway: POST /orders intercepted. Forwarding to gRPC :50052")

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
		}

		var protoReq pbOrders.CreateOrdersRequest
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(bodyBytes, &protoReq)
		if err != nil {
			http.Error(w, "Invalid JSON payload format: "+err.Error(), http.StatusBadRequest)
			return
		}
		res, err := h.ordersClient.CreateOrders(ctx, &protoReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
