package controller

import (
	"context"
	"net/http"
	"restaurantProject/database"
	"restaurantProject/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection = database.OpenConnection("table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		result, err := tableCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		var allTables []bson.M
		err = result.All(ctx, &allTables)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    allTables,
		})

	}
}
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		tableId := c.Param("table_id")

		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "table not found",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"table	": table,
		})

	}
}
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var table models.Table
		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "table not found",
			})

			return
		}
		validateErr := validate.Struct(table)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})

			return
		}
		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.TableId = table.ID.Hex()
		result, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "table not created",
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    result,
		})

	}
}
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var table models.Table
		tableId := c.Param("table_id")
		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "table not found",
			})

			return
		}
		var updateObj primitive.D
		if table.NumberOfGuest != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"number_of_guest": table.NumberOfGuest}})
		}
		if table.TableNumber != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"table_number": table.TableNumber}})
		}
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": table.UpdatedAt}})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		filter := bson.M{"table_id": tableId}
		result, err := tableCollection.UpdateOne(ctx, filter, updateObj, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "table not updated",
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    result,
		})

	}
}
func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
