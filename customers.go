package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("getting customers")
	rows, err := db.Query("SELECT id, name FROM customers order by name")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := []*Customer{}
	for rows.Next() {
		customer := &Customer{}
		err := rows.Scan(&customer.ID, &customer.Name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	result, err := json.Marshal(customers)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("getting customer")
	vars := mux.Vars(r)
	id := vars["customerID"]
	if id == "" {
		log.Println("missing customerID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, name FROM customers WHERE id = $1", id)

	customer := &Customer{}
	err := row.Scan(&customer.ID, &customer.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("adding customer")
	customer := &Customer{}
	err := json.NewDecoder(r.Body).Decode(customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	customer.ID = uuid.New().String()
	_, err = db.Exec("INSERT INTO customers (id, name) VALUES ($1, $2)", customer.ID, customer.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("updating customer")
	customer := &Customer{}
	err := json.NewDecoder(r.Body).Decode(customer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE customers SET name = $1 WHERE id = $2", customer.Name, customer.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("deleting customer")
	vars := mux.Vars(r)
	id := vars["customerID"]
	if id == "" {
		log.Println("missing customerID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAmountOfCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("getting amount of customers")
	row := db.QueryRow("SELECT COUNT(*) FROM customers")

	var amount int
	err := row.Scan(&amount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(amount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
