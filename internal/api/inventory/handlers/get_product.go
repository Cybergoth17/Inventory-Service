package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

// GetProduct retrieves a product by ID.
// @Summary      Get a product by ID
// @Description  Get product details by ID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} inventory.Product
// @Failure      404 {object} inventory.MessageType
// @Failure      500 {object} inventory.MessageType
// @Router       /product/get/{id} [get]
func (h *InventoryHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, product)
}
