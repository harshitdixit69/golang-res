package routes

import (
	"restaurantProject/controller"
	"restaurantProject/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:user_id", middleware.Authentication(), controller.GetUser())
	router.POST("/users/signup", controller.SignUp())
	router.POST("/users/login", controller.Login())
	router.PUT("/users/:user_id", controller.UpdateUser())
}
