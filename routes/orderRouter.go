package routes

import (
	"restaurantProject/controller"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	router.GET("/orders", controller.GetOrders())
	router.GET("/orders/:order_id", controller.GetOrder())
	router.POST("/orders/add", controller.CreateOrder())
	router.PUT("/orders/:order_id", controller.UpdateOrder())
	router.DELETE("/orders/:order_id", controller.DeleteOrder())
}