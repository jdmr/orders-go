package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID       string   `json:"id"`
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
	Price    string   `json:"price"`
}

type Order struct {
	ID       string    `json:"id"`
	Items    []*Item   `json:"items"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
	Customer *Customer `json:"customer"`
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, customer_id, order_date, status FROM orders order by order_date")
	if err != nil {
		log.Printf("Error querying orders: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := &Order{
			Customer: &Customer{},
		}
		err := rows.Scan(&order.ID, &order.Customer.ID, &order.Date, &order.Status)
		if err != nil {
			log.Printf("Error scanning order: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		orders = append(orders, order)
	}

	result, err := json.Marshal(orders)
	if err != nil {
		log.Printf("Error marshalling orders: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["orderID"]
	if id == "" {
		log.Println("Missing order id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing orderID"))
		return
	}

	order := &Order{}
	err := db.QueryRow("SELECT id, customer_id, order_date, status FROM orders where id = $1", id).Scan(&order.ID, &order.Customer.ID, &order.Date, &order.Status)
	if err != nil {
		log.Printf("Error querying order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rows, err := db.Query("SELECT id, product_id, quantity, price FROM order_items where order_id = $1", id)
	if err != nil {
		log.Printf("Error querying order items: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()

	order.Items = []*Item{}
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(&item.ID, &item.Product.ID, &item.Quantity, &item.Price)
		if err != nil {
			log.Printf("Error scanning order item: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		order.Items = append(order.Items, item)
	}

	result, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshalling order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	order := &Order{}
	err := json.NewDecoder(r.Body).Decode(order)
	if err != nil {
		log.Printf("Error decoding order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	order.ID = uuid.New().String()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = tx.QueryRow("INSERT INTO orders (customer_id, order_date, status) VALUES ($1, now(), $2) returning id", order.Customer.ID, order.Status).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for _, item := range order.Items {
		err = tx.QueryRow("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) returning id", order.ID, item.Product.ID, item.Quantity, item.Price).Scan(&item.ID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting order item: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if order.Status == "PAID" {
		amount := ""
		err = tx.QueryRow("SELECT sum(price) FROM order_items WHERE order_id = $1", order.ID).Scan(&amount)
		if err != nil {
			tx.Rollback()
			log.Printf("Error getting order amount: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = tx.Exec("INSERT INTO transactions (order_id, date, amount) VALUES ($1, now(), $2)", order.ID, amount)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshalling order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["orderID"]
	if id == "" {
		log.Println("Missing order id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing orderID"))
		return
	}

	order := &Order{}
	err := json.NewDecoder(r.Body).Decode(order)
	if err != nil {
		log.Printf("Error decoding order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if order.Status == "PAID" {
		log.Println("Cannot update order to PAID")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot update order to PAID"))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = tx.Exec("UPDATE orders SET customer_id = $1, order_date = $2, status = $3 WHERE id = $4", order.Customer.ID, order.Date, order.Status, id)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for _, item := range order.Items {
		_, err = tx.Exec("UPDATE order_items SET product_id = $1, quantity = $2, price = $3 WHERE id = $4", item.Product.ID, item.Quantity, item.Price, item.ID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating order item: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if order.Status == "PAID" {
		amount := ""
		err = tx.QueryRow("SELECT sum(price) FROM order_items WHERE order_id = $1", order.ID).Scan(&amount)
		if err != nil {
			tx.Rollback()
			log.Printf("Error getting order amount: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = tx.Exec("INSERT INTO transactions (order_id, date, amount) VALUES ($1, now(), $2)", order.ID, amount)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting transaction: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if order.Status == "CANCELLED" {
		_, err = tx.Exec("DELETE FROM transactions WHERE order_id = $1", order.ID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			tx.Rollback()
			log.Printf("Error deleting transaction: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshalling order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["orderID"]
	if id == "" {
		log.Println("Missing order id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing orderID"))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = tx.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
}
