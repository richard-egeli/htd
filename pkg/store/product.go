package store

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	SKU         string  `json:"sku"`
	Price       float32 `json:"price"`
}

func (p *Product) Parse(row func(...interface{}) error) error {
	return row(&p.Id, &p.Name, &p.Description, &p.Brand, &p.Category, &p.SKU, &p.Price)
}

func (*Product) Delete(id int) error {
	_, err := store.Exec("DELETE FROM products WHERE id = ?", id)
	return err
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

	return nil
}

func (p *Product) FetchOne(id int) error {
	row := store.QueryRow("SELECT * FROM products WHERE id = ?", id)
	return p.Parse(row.Scan)
}

func (*Product) FetchAll() ([]Product, error) {
	var products []Product
	rows, err := store.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product Product
		err := product.Parse(rows.Scan)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (*Product) FetchPage(size, page int) ([]Product, error) {
	var products []Product

	rows, err := store.Query("SELECT * FROM products LIMIT ? OFFSET ?", size, size*page)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product Product
		err := product.Parse(rows.Scan)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
