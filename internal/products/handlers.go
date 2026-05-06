package products

import (
	"log"
	"net/http"
	"strconv"
	"time"

	repo "github.com/Prakash-Ravichandran/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *Handler) ListProductsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// json.Write(w, http.StatusOK, id)
	idnum, errnum := strconv.ParseInt(id, 10, 64)
	if errnum != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	product, err := h.service.ListProductsByID(r.Context(), idnum)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)

}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	products := repo.CreateProductParams{ID: 17, Name: "Omen", PriceInCents: 55, Quantity: 10, CreatedAt: pgtype.Timestamptz{
		Time:  time.Now(), // This is the standard way
		Valid: true,
	}}
	productMessage, err := h.service.CreateProducts(r.Context(), products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.Write(w, http.StatusOK, productMessage)
}
