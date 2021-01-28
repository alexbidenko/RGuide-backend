package controllers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"rguide/dif"
	"rguide/models"
	_ "rguide/models"
)

func InitProducts(r *mux.Router) {
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
}

func getProducts(w http.ResponseWriter, _ *http.Request) {
	var product []models.Product
	dif.DB.Select(&product, "SELECT * FROM products")
	data, _ := json.Marshal(product)
	w.Write(data)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = dif.DB.QueryRowx(`INSERT INTO products (title, description) VALUES ($1, $2) RETURNING id`, p.Title, p.Description).StructScan(&p)
	data, _ := json.Marshal(p)
	w.Write(data)
}
