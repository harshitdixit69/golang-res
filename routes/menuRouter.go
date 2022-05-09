package routes

import (
	"restaurantProject/controller"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	router.GET("/menus", controller.GetMenus())
	router.GET("/menus/:menu_id", controller.GetMenu())
	router.POST("/menus/add", controller.CreateMenu())
	router.PUT("/menus/:menu_id", controller.UpdateMenu())
	router.DELETE("/menus/:menu_id", controller.DeleteMenu())
}
