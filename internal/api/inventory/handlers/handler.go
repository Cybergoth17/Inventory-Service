package handlers

import inventoryclient "inventory-service/internal/service/inventory"

type InventoryHandler struct {
	service inventoryclient.Client
}

func NewInventoryHandler(service inventoryclient.Client) *InventoryHandler {
	return &InventoryHandler{service: service}
}
