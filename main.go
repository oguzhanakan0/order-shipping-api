package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oguzhanakan0/order-shipping-api/api"
)

// Returns a `gin.Engine` router including the predefined endpoints and
// loads the static files with html templates
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.tmpl")
	r.Static("/assets", "./assets")

	// HTML Responses
	r.GET("/", api.IndexHTML)
	r.GET("/order", api.GetShipmentHTML)

	// JSON Responses
	r.GET("/api/order", api.GetShipment)

	return r
}

// Runs the API service
func main() {
	r := setupRouter()
	r.Run(":8080")
}
