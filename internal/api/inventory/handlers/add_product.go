package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"inventory-service/internal/model/inventory"
)

// AddProduct adds a new product to the inventory.
// @Summary      Add a new product
// @Description  Add a new product to the inventory
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param       product body inventory.Product true "Product object"
// @Success     201 {object} inventory.Product
// @Failure     400 {object} inventory.MessageType
// @Failure     500 {object} inventory.MessageType
// @Router      /product/add [post]
func (h *InventoryHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product inventory.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, product)
}
