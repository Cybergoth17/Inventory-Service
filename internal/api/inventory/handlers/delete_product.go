package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

// DeleteProduct removes a product by ID.
// @Summary      Delete a product by ID
// @Description  Delete product from the inventory by ID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      204
// @Failure      500 {object} inventory.MessageType
// @Router       /product/delete/{id} [delete]
func (h *InventoryHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusNoContent)
}
