package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"rguide/entities"
	_ "rguide/entities"
	"rguide/models"
	"strconv"
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

	parameters := map[string]interface{}{}

	categoryIdParameter := r.URL.Query().Get("category_id")
	categoryId, err := strconv.Atoi(categoryIdParameter)

	if categoryIdParameter != "" && err == nil {
		parameters["category_id"] = categoryId
	}

	var productModel models.ProductModel
	products := productModel.FindAll(q, parameters)

	data, _ := json.Marshal(products)
	w.Write(data)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var productModel models.ProductModel
	product, err := productModel.GetById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(product)
	w.Write(data)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p entities.Product
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

	var productModel models.ProductModel
	productModel.Create(&p)

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

	productId, err := strconv.Atoi(r.FormValue("productId"))
	if err == nil {
		var productModel models.ProductModel
		product, _ := productModel.GetById(productId)

		var prevFileName string
		switch fileType {
		case "model":
			prevFileName = product.Model
			product.Model = fileName
			break
		case "preview":
			prevFileName = product.Preview
			product.Preview = fileName
		}
		os.Remove("/var/www/files-" + fileType + "/" + prevFileName)

		productModel.Update(productId, &product)
	}

	fmt.Fprintf(w, `{"filename":"` + fileName + `"}`)
}
