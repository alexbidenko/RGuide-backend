package controllers

import (
	"encoding/json"
	"net/http"
	"rguide/dif"
	"rguide/models"
	_ "rguide/models"
)

func GetProducts(w http.ResponseWriter, _ *http.Request) {
	var product []models.Product
	dif.DB.Select(&product, "SELECT * FROM products")
	data, _ := json.Marshal(product)
	w.Write(data)
}
