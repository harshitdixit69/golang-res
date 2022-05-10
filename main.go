package main

import (
	"restaurantProject/middleware"
	"restaurantProject/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	router.Run(":" + port)
}
