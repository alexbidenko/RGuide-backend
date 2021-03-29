package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/doug-martin/goqu"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"rguide/dif"
	"rguide/models"
	_ "rguide/models"
	"strings"
)

func InitProducts(r *mux.Router) {
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}", getProductById).Methods("GET")
	r.HandleFunc("/products/file", uploadFile).Methods("POST")
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	products := make([]models.Product, 0)
	builder := goqu.From("products")
	var query string
	if q != "" {
		builder = builder.Where(goqu.Ex{
			"title": goqu.Op{"like": "%" + q + "%"},
		})
	}
	query, _, _ = builder.ToSQL()
	dif.DB.Select(&products, query)
	data, _ := json.Marshal(products)
	w.Write(data)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	sql, _, _ := goqu.From("products").Where(goqu.Ex{"id": id}).ToSQL()

	var product models.Product
	dif.DB.Get(&product, sql)

	if product.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	sql, _, _ := goqu.Insert("products").Rows(p).ToSQL()

	_ = dif.DB.QueryRowx(sql).StructScan(&p)
	data, _ := json.Marshal(p)
	w.Write(data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 21)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileType := r.FormValue("type")
	if fileType != "preview" && fileType != "model" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tempFile, err := ioutil.TempFile("/var/www/files-" + fileType, "file-*-" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)
	fileName := strings.ReplaceAll(tempFile.Name(), "/var/www/files-" + fileType + "/", "")

	productId := r.FormValue("productId")
	if productId != "" {
		sql, _, _ := goqu.From("products").Where(goqu.Ex{"id": productId}).ToSQL()
		var product models.Product
		dif.DB.Get(&product, sql)

		var prevFileName string
		switch fileType {
		case "model":
			prevFileName = product.Model
			break
		case "preview":
			prevFileName = product.Preview
		}
		os.Remove("/var/www/files-" + fileType + "/" + prevFileName)

		sql, _, _ = goqu.Update("products").Set(goqu.Record{fileType: fileName}).Where(goqu.Ex{"id": productId}).ToSQL()
		dif.DB.Exec(sql)
	}

	fmt.Fprintf(w, `{"filename":"` + fileName + `"}`)
}
