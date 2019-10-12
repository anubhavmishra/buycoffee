package inventory

import (
	"encoding/json"
	"net/http"
)

type InventoryResponse struct {
	Inventory []InventoryItem `json:"inventory"`
}

type InventoryItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ImageURL string  `json:"image_url"`
	Price    float32 `json:"price"`
	Currency string  `json:"currency"`
}

var inventoryAvailable = []InventoryItem{
	{ID: 1, Name: "Flat White", ImageURL: "/static/images/flat-white.png", Currency: "USD", Price: float32(5.50)},
	{ID: 2, Name: "Black Coffee", ImageURL: "/static/images/black-coffee.png", Currency: "USD", Price: float32(2.5)},
	{ID: 3, Name: "Americano", ImageURL: "/static/images/americano.jpg", Currency: "USD", Price: float32(3.50)},
	{ID: 4, Name: "Latte", ImageURL: "/static/images/latte.jpg", Currency: "USD", Price: float32(5.80)},
	{ID: 5, Name: "Cappucino", ImageURL: "/static/images/cappuccino.jpg", Currency: "USD", Price: float32(5.50)},
}

func InventoryHandler(w http.ResponseWriter, r *http.Request) {

	response := InventoryResponse{
		Inventory: inventoryAvailable,
	}
	json.NewEncoder(w).Encode(response)
	return
}
