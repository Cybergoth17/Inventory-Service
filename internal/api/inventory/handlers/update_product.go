package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"inventory-service/internal/model/inventory"
	"net/http"
)

// UpdateProduct updates a product by ID.
// @Summary      Update a product by ID
// @Description  Update product details by ID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Param        updatedProduct body inventory.Product true "Updated product object"
// @Success      200 {object} inventory.Product
// @Failure      400 {object} inventory.MessageType
// @Failure      500 {object} inventory.MessageType
// @Router       /product/update/{id} [put]
func (h *InventoryHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updatedProduct inventory.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, &updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedProduct)
}
