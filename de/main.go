package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"hcmut.vn/de/handlers"
	"hcmut.vn/de/resources"
)

func main() {
	// Load environment config
	cfg := newConfig()

	// // Connect to oracle server
	db, err := resources.NewConnection(cfg.oracleDataSourceName)
	if err != nil {
		log.Fatal("Connect to oracle error: ", err.Error())
	}
	defer db.Close()

	ph := handlers.NewProductHandler(db)
	ch := handlers.NewCartHandler(db)

	r := gin.Default()

	r.GET("/api/product/:product_id", ph.GetProduct)
	r.POST("/api/product", ph.CreateProduct)
	r.POST("/api/cart", ch.AddToCart)

	_ = r.Run(":" + cfg.ginPort)
}
