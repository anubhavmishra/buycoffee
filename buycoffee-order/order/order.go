package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type OrderResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	OrderID int    `json:"order_id"`
}

type orderHandler struct {
	db *sql.DB
}

func (o *orderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response := OrderResponse{}

	switch r.Method {
	case "GET":
		response.Error = "Currently unsupported."
		w.WriteHeader(http.StatusBadRequest)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			response.Error = fmt.Sprintf("ParseForm() err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			break
		}
		userID := r.FormValue("user_id")
		itemName := r.FormValue("item")

		fmt.Println("user_id:", userID, "item_name:", itemName)

		if userID == "" || itemName == "" {
			response.Error = "invalid form request."
			w.WriteHeader(http.StatusBadRequest)
			break
		}

		// convert user_id value before db insert
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			response.Error = "form validation for user id failed."
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		orderID, err := o.placeOrder(userIDInt, itemName)
		if err != nil {
			response.Error = "Something went wrong!"
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		response.Message = "order successfully created!"
		response.OrderID = orderID

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response.Error = "Sorry, only GET and POST methods are supported"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func OrderHandler(db *sql.DB) http.Handler {
	return &orderHandler{
		db: db,
	}
}

func (o *orderHandler) placeOrder(userID int, itemName string) (int, error) {
	var orderID int
	err := o.db.QueryRow("INSERT INTO \"order\"(user_id, item_name, created_on) VALUES ($1,$2, NOW()) RETURNING order_id", userID, itemName).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
