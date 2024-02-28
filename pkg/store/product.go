package store

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	SKU         string `json:"sku"`
	Price       int    `json:"price"`
}

func (p *Product) Insert() error {
	stmt, err := store.Prepare("INSERT INTO products(name, description, brand, category, sku, price) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(p.Name, p.Description, p.Brand, p.Category, p.SKU, p.Price)
	if err != nil {
		return err
	}

	fmt.Println("Successfully added")
	return nil
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received product: %+v\n", product)

	if err := product.Insert(); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product added successfully"))
}

func FetchProductsAll(w http.ResponseWriter, r *http.Request) {
	var products []Product
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	rows, err := store.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Brand, &p.Category, &p.SKU, &p.Price)
		if err != nil {
			http.Error(w, "Issue with one of the found rows", http.StatusInternalServerError)
			return
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Problem with internal rows", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(products)
}
