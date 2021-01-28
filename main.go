package main

import (
	_ "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rguide/controllers"
	"rguide/dif"
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

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	controllers.InitProducts(s)
	r.HandleFunc("/api", handler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8010", nil))
}