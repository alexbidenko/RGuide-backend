package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rguide/dif"
	"rguide/models"
	_ "rguide/models"
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
	if q != "" {
		dif.DB.Select(&products, "SELECT * FROM products WHERE title LIKE $1", "%" + q + "%")
	} else {
		dif.DB.Select(&products, "SELECT * FROM products")
	}
	data, _ := json.Marshal(products)
	w.Write(data)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var product models.Product
	dif.DB.Get(&product, "SELECT * FROM products WHERE id = $1", id)

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

	_ = dif.DB.QueryRowx(`INSERT INTO products (title, description) VALUES ($1, $2) RETURNING id`, p.Title, p.Description).StructScan(&p)
	data, _ := json.Marshal(p)
	w.Write(data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 21)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	if r.FormValue("type") != "preview" && r.FormValue("type") != "model" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("/var/www/files-" + r.FormValue("type"), "file-*-" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, `{"filename":"` + tempFile.Name() + `"}`)
}
