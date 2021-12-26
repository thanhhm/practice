package handlers

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"hcmut.vn/de/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{
		productService: services.NewProductService(db),
	}
}

func (pdh *ProductHandler) CreateProduct(c *gin.Context) {
	var reqBody services.CreateProductRequest
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, "Parsing request error")
		return
	}
	if reqBody.ProductID <= 0 {
		c.JSON(400, "Invalid product_id")
		return
	}

	err = pdh.productService.CreateProduct(reqBody)
	if err != nil {
		log.Println("Create product error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func (pdh *ProductHandler) GetProduct(c *gin.Context) {
	sortCondition := c.Query("sort")
	req := services.GetProductsRequest{
		SortCondition: sortCondition,
	}

	result, err := pdh.productService.GetProduct(req)
	if err != nil {
		log.Println("Create product error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, result)
}
