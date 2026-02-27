package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/handler"
	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, products!")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	storage.InitDB(conn)
	router := gin.Default()
	router.POST("/products", handler.CreateProduct)
	router.PUT("/products/:id", handler.UpdateProduct)
	router.DELETE("/products/:id", handler.DeleteProduct)
	router.GET("/products", handler.GetAllProducts)
	router.Run(":8080")

}
