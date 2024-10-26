package api

import (
	"inventory-service/internal/api/inventory/handlers"
	"net/http"
)

type API struct {
	inventoryHandler *handlers.InventoryHandler
}

func NewAPI(inventoryHandler *handlers.InventoryHandler) *API {
	return &API{inventoryHandler: inventoryHandler}
}

func (api *API) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/product/add", api.inventoryHandler.AddProduct)
	mux.HandleFunc("/product/update", api.inventoryHandler.UpdateProduct)
	mux.HandleFunc("/product/delete", api.inventoryHandler.DeleteProduct)
	mux.HandleFunc("/product/get", api.inventoryHandler.GetProduct)

	return mux
}
