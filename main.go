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

var headers = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header"

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

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", headers)
		w.WriteHeader(http.StatusNoContent)
		return
	})
	r.Use(accessControlMiddleware)

	r.PathPrefix("/files/previews/").Handler(http.StripPrefix("/files/previews/", http.FileServer(http.Dir("/var/www/files-preview"))))
	r.PathPrefix("/files/models/").Handler(http.StripPrefix("/files/models/", http.FileServer(http.Dir("/var/www/files-model"))))
	s := r.PathPrefix("/api").Subrouter()
	controllers.InitProducts(s)
	r.HandleFunc("/api", handler)
	fmt.Printf("Server starting")
	log.Fatal(http.ListenAndServe(":8010", r))
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", headers)
		log.Printf(r.Method)

		next.ServeHTTP(w, r)
	})
}
