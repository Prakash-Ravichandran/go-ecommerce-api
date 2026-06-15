package main

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/orders"
	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/products"
	pb "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/health"
	pbOrders "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/orders"
	pbProduct "github.com/Prakash-Ravichandran/go-ecommerce-api/proto/product"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type application struct {
	config config
	db     *pgx.Conn
	// mount
	// run
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting - random ID for each request
	r.Use(middleware.RealIP)    // get the real IP address of the client
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //  recover from crashes

	// set a timeout value on the request context(ctx), that will signal
	// through ctx.Done() that the request has timed out
	// and further processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second)) // if request takes more than 60s then stop it

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productsService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productsService) // pass the service
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.ListProductsByID)
	r.Post("/products", productHandler.HandleCreateProduct)
	r.Put("/products", productHandler.HandleUpdateProduct)

	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	r.Get("/orders", ordersHandler.HandleGetOrders)
	r.Get("/orders/{id}", ordersHandler.HandleFindOrder)
	r.Post("/orders", ordersHandler.HandlePostOrders)
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("Starting server on", "addr", app.config.addr)

	return srv.ListenAndServe()
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// GRPC Products server
func (app *application) runProductsGRPC() error {
	grpcAddr := ":50051"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()

	productRepo := repo.New(app.db)
	productHandler := &grpcServer{
		db:   app.db,
		repo: *productRepo,
	}

	pbProduct.RegisterProductServiceServer(gServer, productHandler)
	slog.Info("Starting gRPC Products server on", "addr", grpcAddr)

	return gServer.Serve(lis)
}

// GRPC Orders server
func (app *application) runOrderGRPC() error {
	orderGRPCAddr := ":50052"

	lis, err := net.Listen("tcp", orderGRPCAddr)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()

	orderRepo := repo.New(app.db)
	ordersHandler := &grpcServer{
		db:   app.db,
		repo: *orderRepo,
	}

	slog.Info("Registering Orders Service Server")
	pbOrders.RegisterOrderServiceServer(gServer, ordersHandler)

	slog.Info("Starting gRPC Orders server on", "addr", orderGRPCAddr)

	return gServer.Serve(lis)
}

// GRPC Health server
func (app *application) runHealthGRPC() error {
	healthGRPCAddr := ":50053"
	lis, err := net.Listen("tcp", healthGRPCAddr)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()

	slog.Info("Registering Health Service Server")
	pb.RegisterHealthServiceServer(gServer, &grpcServer{})
	slog.Info("Starting gRPC Health server on", "addr", healthGRPCAddr)
	return gServer.Serve(lis)
}
