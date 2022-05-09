package controller

import (
	"context"
	"log"
	"net/http"
	"restaurantProject/database"
	"restaurantProject/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection = database.OpenConnection("menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		result, err := menuCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occur while listing items",
			})
		}
		var allMenus []bson.M
		if err = result.All(ctx, allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"message": result,
		})
	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		menuId := c.Param("menu_id")

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"food_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "menu not found",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"menu": menu,
		})

	}
}
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()

		var menu models.Menu
		bindJsonErr := c.BindJSON(&menu)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})
			return
		}
		validateErr := validate.Struct(menu)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()
		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": result,
		})

	}
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var menu models.Menu
		bindJsonErr := c.BindJSON(&menu)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})

			return
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}
		var updateObj primitive.D
		if menu.StartDate != nil && menu.EndDate != nil {
			if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Start date must be less than end date",
				})

				return
			}
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"start_date": menu.StartDate, "end_date": menu.EndDate}})
			if menu.Name != "" {
				updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"name": menu.Name}})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"category": menu.Category}})
			}
			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": menu.UpdatedAt}})
			upsertL := true
			opt := options.UpdateOptions{
				Upsert: &upsertL,
			}
			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{Key: "$set", Value: updateObj},
				},
				&opt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": result,
			})
		}

	}
}
func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}
