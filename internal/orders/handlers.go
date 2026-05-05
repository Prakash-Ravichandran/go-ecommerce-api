package orders

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/json"
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
	OrdersFromDb, err := h.service.GetOrders(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, OrdersFromDb)
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
