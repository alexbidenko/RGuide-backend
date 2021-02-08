package main

import (
	_ "database/sql"
	"fmt"
	"github.com/gorilla/handlers"
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

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	r := mux.NewRouter()
	r.PathPrefix("/files/previews").Handler(http.FileServer(http.Dir("/var/www/preview-images")))
	r.PathPrefix("/files/models").Handler(http.FileServer(http.Dir("/var/www/model-files")))
	s := r.PathPrefix("/api").Subrouter()
	controllers.InitProducts(s)
	r.HandleFunc("/api", handler)
	http.Handle("/", r)
	fmt.Printf("Server starting")
	log.Fatal(http.ListenAndServe(":8010", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}