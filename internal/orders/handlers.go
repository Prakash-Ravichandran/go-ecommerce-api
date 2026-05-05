package orders

import (
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service OrderService
}

func NewHandler(o OrderService) *handler {
	return &handler{
		service: o,
	}
}

func (h *handler) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	OrdersFromDb, err := h.service.ListOrders(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, OrdersFromDb)
}

func (h *handler) HandleFindOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idnum, errnum := strconv.ParseInt(id, 10, 64)
	if errnum != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	order, err := h.service.FindOrderById(r.Context(), idnum)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, order)
}

func (h *handler) HandlePostOrders(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams

	if err := json.Read(r, &tempOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("service error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)

}
