package repository

import (
	"database/sql"
	"github.com/lucasreisprestes/application-go-hexagonal/internal/entity"
)

type ProductRepositoryMysql struct {
	DB *sql.DB
}

// function
func NewProductRepositoryMysql(db *sql.DB) *ProductRepositoryMysql {
	return &ProductRepositoryMysql{DB: db}
}

// method
func (r *ProductRepositoryMysql) Create(product *entity.Product) error {
	_, err := r.DB.Exec("Insert into products (id, name, price) values(?,?,?)",
		product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepositoryMysql) FindAll() ([]*entity.Product, error) {
	rows, err := r.DB.Query("select id, name, price from products")

	if err != nil {
		return nil, err
	}

	//espera executar tudo, para depois fechar a conexão
	defer rows.Close()
	var products []*entity.Product

	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)

		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
