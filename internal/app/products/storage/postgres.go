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
