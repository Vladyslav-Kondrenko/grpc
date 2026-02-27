package storage

import (
	"context"

	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/product"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func InitDB(conn *pgx.Conn) {
	db = conn
}

func GetAllProducts(ctx context.Context) ([]product.Product, error) {
	rows, err := db.Query(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []product.Product{}
	for rows.Next() {
		var product product.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func CreateProduct(ctx context.Context, p product.ProductCreateRequest) (product.Product, error) {
	var createdProduct product.Product
	err := db.QueryRow(ctx,
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id, name, price",
		p.Name, p.Price,
	).Scan(&createdProduct.ID, &createdProduct.Name, &createdProduct.Price)
	if err != nil {
		return createdProduct, err
	}
	return createdProduct, nil
}

func UpdateProduct(ctx context.Context, p product.ProductUpdateRequest) (product.Product, error) {
	var updatedProduct product.Product
	err := db.QueryRow(ctx,
		"UPDATE products SET name = $1, price = $2 WHERE id = $3 RETURNING id, name, price",
		p.Name, p.Price, p.ID,
	).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Price)
	if err != nil {
		return updatedProduct, err
	}
	return updatedProduct, nil
}

func DeleteProduct(ctx context.Context, id int) error {
	_, err := db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
