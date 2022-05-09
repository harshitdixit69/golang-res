package routes

import (
	"restaurantProject/controller"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.Engine) {
	router.GET("/foods", controller.GetFoods())
	router.POST("/foods/add", controller.CreateFood())
	router.GET("/foods/:food_id", controller.GetFood())
	router.PUT("/foods/:food_id", controller.UpdateFood())
	router.DELETE("/foods/:food_id", controller.DeleteFood())
}
