package routes

import (
	"restaurantProject/controller"

	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.Engine) {
	router.GET("/tables", controller.GetTables())
	router.GET("/tables/:table_id", controller.GetTable())
	router.POST("/tables/add", controller.CreateTable())
	router.PUT("/tables/:table_id", controller.UpdateTable())
	router.DELETE("/tables/:table_id", controller.DeleteTable())
}
