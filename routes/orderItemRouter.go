package routes

import (
	"restaurantProject/controller"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(router *gin.Engine) {
	router.GET("/orderItems", controller.GetOrderItems())
	router.GET("/orderItems/:order_item_id", controller.GetOrderItem())
	router.GET("/orderItems-order/:order_id", controller.GetOrderItemsByOrder())
	router.POST("/orderItems/add", controller.CreateOrderItem())
	router.PUT("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}
