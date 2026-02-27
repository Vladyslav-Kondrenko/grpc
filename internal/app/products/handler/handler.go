package handler

import (
	"net/http"

	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/storage"
	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	products, err := storage.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	panic("CreateProduct is not implemented")
}

func UpdateProduct(c *gin.Context) {
	panic("UpdateProduct is not implemented")
}

func DeleteProduct(c *gin.Context) {
	panic("DeleteProduct is not implemented")
}
