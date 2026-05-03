package orders

import (
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
	dummyOrdersFromDb := h.service.GetOrder(r.Context())
	json.Write(w, http.StatusOK, dummyOrdersFromDb)
}
