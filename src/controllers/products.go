package controllers

import (
	"dif"
	"encoding/json"
	"models"
	_ "models"
	"net/http"
)

func GetProducts(w http.ResponseWriter, _ *http.Request) {
	var product []models.Product
	dif.DB.Select(&product, "SELECT * FROM products")
	data, _ := json.Marshal(product)
	w.Write(data)
}
