package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type ImageType string

type Image struct {
	Id   int
	Name string
	Alt  string
	Type string
	Path string
}

func generateHash() string {
	randomBytes := make([]byte, 5)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating hash:", err)
		return ""
	}

	return hex.EncodeToString(randomBytes) + "."
}

func (image *Image) Insert() error {
	stmt, err := store.Prepare("INSERT INTO images(name, alt, type, path) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(image.Name, image.Alt, image.Type, image.Path)
	if err != nil {
		return err
	}

	return nil
}

func (image *Image) Fetch(id int) error {
	row := store.QueryRow("SELECT * FROM images WHERE id = ?", id)
	return row.Scan(&image.Id, &image.Name, &image.Alt, &image.Type, &image.Path)
}

func (image *Image) Delete() error {
	_, err := store.Exec("DELETE FROM images WHERE id = ?", image.Id)
	return err
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

	fmt.Println("ID", id)
	fmt.Fprintf(w, "Image deleted")
}
