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
	"strings"
)

var headers = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, sec-ch-ua, sec-ch-ua-mobile, sec-fetch-dest, sec-fetch-mode, sec-fetch-site"

func handler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello RGuide!")
}

func main() {
	err := dif.DBError
	if err != nil {
		panic(err)
	}

	dif.Migrate()

	r := mux.NewRouter()

	headersOk := handlers.AllowedHeaders(strings.Split(headers, ", "))
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "PATCH"})

	r.PathPrefix("/files/previews/").Handler(http.StripPrefix("/files/previews/", http.FileServer(http.Dir("/var/www/files-preview"))))
	r.PathPrefix("/files/models/").Handler(http.StripPrefix("/files/models/", http.FileServer(http.Dir("/var/www/files-model"))))

	s := r.PathPrefix("/api").Subrouter()
	controllers.InitProducts(s)
	controllers.InitCategories(s)
	r.HandleFunc("/api", handler)

	fmt.Printf("Server starting")
	log.Fatal(http.ListenAndServe(":8010", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
