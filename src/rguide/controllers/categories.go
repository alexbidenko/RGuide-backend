package controllers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"rguide/entities"
	"rguide/models"
)

func InitCategories(r *mux.Router) {
	r.HandleFunc("/categories", getCategories).Methods("GET")
	r.HandleFunc("/categories", createCategory).Methods("POST")
}

func getCategories(w http.ResponseWriter, _ *http.Request) {
	var categoryModel models.CategoryModel
	categories := categoryModel.All()

	data, _ := json.Marshal(categories)
	w.Write(data)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var c entities.Category
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var categoryModel models.CategoryModel
	categoryModel.Create(&c)

	data, _ := json.Marshal(c)
	w.Write(data)
}
