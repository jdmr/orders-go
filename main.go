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

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
