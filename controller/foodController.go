package controller

import (
	"context"
	"log"
	"math"
	"net/http"
	"restaurantProject/database"
	"restaurantProject/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection = database.OpenConnection("food")

var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx = context.Background()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		matchStage := primitive.D{
			primitive.E{Key: "$match", Value: bson.M{}},
		}
		groupStage := primitive.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"_id": "null",
				},
				"total_count": bson.M{
					"$sum": 1,
				},
				"data": bson.M{
					"$push": "$$ROOT",
				},
			}},
		}
		projectStage := bson.D{
			primitive.E{Key: "$project", Value: bson.M{
				"_id":         0,
				"total_count": 1,
				"totalcount":  1,
				"food_items":  "$data",
			}},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			return
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		foodId := c.Param("food_id")

		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "food not found",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"food": food,
		})

	}
}
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()

		var menu models.Menu
		var food models.Food
		bindJsonErr := c.BindJSON(&food)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})
			return
		}
		validateErr := validate.Struct(food)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}
		foodId := *food.MenuId
		err := menuCollection.FindOne(ctx, bson.M{"menuid": foodId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "menu not found",
			})
			return
		}
		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Food Not inserted",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
func round(num float64) int {
	return int(num + 0.5)
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var menu models.Menu
		var food models.Food
		foodId := c.Param("food_id")
		bindJsonErr := c.BindJSON(&food)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})

			return
		}
		var updateObj primitive.D
		if food.Name != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"name": food.Name}})
		}
		if food.Price != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"price": food.Price}})
		}
		if food.FoodImage != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"food_image": food.FoodImage}})
		}
		if food.MenuId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "menu not found",
				})

				return
			}
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"menu_id": food.MenuId}})
		}
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": food.UpdatedAt}})

		upsertL := true
		filter := bson.M{"food_id": foodId}
		opt := options.UpdateOptions{
			Upsert: &upsertL,
		}
		result, err := foodCollection.UpdateOne(ctx, filter, updateObj, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Food Not updated",
			})

			return
		}
		c.JSON(http.StatusOK, result)

	}
}
func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
