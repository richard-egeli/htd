package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

func (image *Image) Parse(row func(...interface{}) error) error {
	return row(&image.Id, &image.Name, &image.Alt, &image.Type, &image.Path)
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
	return image.Parse(row.Scan)
}

func (image *Image) Delete() error {
	_, err := store.Exec("DELETE FROM images WHERE id = ?", image.Id)
	return err
}
