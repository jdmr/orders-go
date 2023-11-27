package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Connecting to database")
	var err error
	db, err = sql.Open("pgx", "postgres://jdmr:test@localhost:5432/orders?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating new router")
	router := mux.NewRouter()

	log.Println("Registering routes")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{customerID}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{customerID}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{customerID}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customer-amount", getAmountOfCustomers).Methods("GET")

	router.HandleFunc("/products", getProducts).Methods("GET")
	router.HandleFunc("/products", addProduct).Methods("POST")
	router.HandleFunc("/products/{productID}", getProduct).Methods("GET")
	router.HandleFunc("/products/{productID}", updateProduct).Methods("PUT")
	router.HandleFunc("/products/{productID}", deleteProduct).Methods("DELETE")

	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/orders", addOrder).Methods("POST")
	router.HandleFunc("/orders/{orderID}", getOrder).Methods("GET")
	router.HandleFunc("/orders/{orderID}", updateOrder).Methods("PUT")
	router.HandleFunc("/orders/{orderID}", deleteOrder).Methods("DELETE")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
