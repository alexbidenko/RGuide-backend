package main

import (
	"controllers"
	_ "database/sql"
	"dif"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"migrations"
	"models"
	"net/http"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello RGuide!")
}

func main() {
	db := dif.DB
	err := dif.DBError
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec(migrations.Schema)
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO products (id, title, description) VALUES (:id, :title, :description)", &models.Product{Id: 1, Title: "title", Description: "description"})
	tx.Commit()

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/products", controllers.GetProducts)
	r.HandleFunc("/api", handler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}