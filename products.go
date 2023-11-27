package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products order by name")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		products = append(products, product)
	}

	result, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing productID"))
		return
	}

	row := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id)

	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	product := &Product{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if product.ID == "" || product.Name == "" || product.Price == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing productID, name, or price"))
		return
	}

	_, err = db.Exec("INSERT INTO products (id, name, price) VALUES ($1, $2, $3)", product.ID, product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing productID"))
		return
	}

	product := &Product{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	product.ID = id

	if product.ID == "" || product.Name == "" || product.Price == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing productID, name, or price"))
		return
	}

	_, err = db.Exec("UPDATE products SET name = $2, price = $3 WHERE id = $1", product.ID, product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing productID"))
		return
	}

	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
