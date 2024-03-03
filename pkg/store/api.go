package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)

	if err != nil {
		http.Error(w, "Invalid price tag", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	product.Name = r.FormValue("name")
	product.Brand = r.FormValue("brand")
	product.Description = r.FormValue("description")
	product.Category = r.FormValue("category")
	product.SKU = r.FormValue("sku")
	product.Price = float32(price)

	if err := product.Insert(); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Product added successfully")
}

func FetchProducts(w http.ResponseWriter, r *http.Request) {
	var handler Product
	var products []Product

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	if page == "" && size == "" {
		p, err := handler.FetchAll()
		if err != nil {
			http.Error(w, "Failed to get products", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		products = append(products, p...)
	} else {
		page, pageErr := strconv.Atoi(page)
		size, sizeErr := strconv.Atoi(size)

		if pageErr != nil || sizeErr != nil {
			http.Error(w, "Invalid query parameters", http.StatusBadGateway)
			fmt.Println(pageErr, sizeErr)
			return
		}

		p, err := handler.FetchPage(size, page)
		if err != nil {
			http.Error(w, "Failed to get page", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		products = append(products, p...)
	}

	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(products)

	if err != nil {
		http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
	}
}

func FetchProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid path id", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = product.FetchOne(id)
	if err != nil {
		http.Error(w, "Error fetching product", http.StatusNotFound)
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	var image Image
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		return
	}

	err = image.Fetch(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "No image found with id", http.StatusNotFound)
		return
	}

	err = image.Delete()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to remove image", http.StatusInternalServerError)
		return
	}

	err = os.Remove(image.Path)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to remove the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Image deleted")
}

func CreateImage(w http.ResponseWriter, r *http.Request) {
	var image Image

	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	image.Name = r.FormValue("name")
	image.Alt = r.FormValue("alt")
	image.Type = filepath.Ext(fileHeader.Filename)
	image.Path = string(http.Dir("./static/" + generateHash() + fileHeader.Filename))

	osFile, err := os.Create(image.Path)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	defer osFile.Close()

	_, err = osFile.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	if err := image.Insert(); err != nil {
		fmt.Println(err)
		http.Error(w, "Error storing image", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Image uploaded successfully")
}
