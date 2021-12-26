package handlers

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"hcmut.vn/de/services"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(db *sql.DB) *CartHandler {
	return &CartHandler{
		cartService: services.NewCartService(db),
	}
}

func (ch *CartHandler) AddToCart(c *gin.Context) {
	var reqBody services.AddToCartRequest
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, "Parsing request error")
		return
	}
	if reqBody.CartID <= 0 {
		c.JSON(400, "Invalid cart_id")
		return
	}

	err = ch.cartService.AddToCart(reqBody)
	if err != nil {
		log.Println("Add to cart error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}
