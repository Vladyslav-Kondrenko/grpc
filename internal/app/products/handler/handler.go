package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/product"
	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/products/storage"
	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	writer := csv.NewWriter(c.Writer)
	writer.Comma = ';'

	products, err := storage.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=products.csv")
	writer.Write([]string{"PRODUCT NAME", "PRICE"})
	for _, p := range products {
		writer.Write([]string{p.Name, strconv.Itoa(p.Price)})
	}

	writer.Flush()
}

func CreateProduct(c *gin.Context) {
	var p product.ProductCreateRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := storage.CreateProduct(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var p product.ProductUpdateRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = idInt
	product, err := storage.UpdateProduct(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = storage.DeleteProduct(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
