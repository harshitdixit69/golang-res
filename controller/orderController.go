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

var orderCollection = database.OpenConnection("order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		result, err := orderCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		var allOrder []bson.M
		err = result.All(ctx, &allOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    allOrder,
		})

	}
}
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		orderId := c.Param("order_id")

		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "food not found",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"order	": order,
		})

	}
}
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()

		var table models.Table
		var order models.Order
		bindJsonErr := c.BindJSON(&order)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})
			return
		}
		validateErr := validate.Struct(order)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}
		if order.TableId != nil {
			err := tableCollection.FindOne(ctx, bson.M{"tableid": order.TableId}).Decode(&table)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "table not found",
				})
				return
			}
		}
		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()
		insertResult, err := orderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    insertResult.InsertedID,
		})
	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var table models.Table
		var order models.Order
		var updateObj primitive.D
		orderId := c.Param("order_id")
		bindJsonErr := c.BindJSON(&order)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})

			return
		}
		if order.TableId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "menu not found",
				})

				return
			}
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"menu_id": order.TableId}})
			order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": order.UpdatedAt}})
			upsertL := true
			filter := bson.M{"order_id": orderId}
			opt := options.UpdateOptions{
				Upsert: &upsertL,
			}
			result, err := orderCollection.UpdateOne(ctx, filter, updateObj, &opt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Food Not updated",
				})

				return
			}
			c.JSON(http.StatusOK, result)

		}

	}
}
func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
func orderItemOrderCreator(order models.Order) string {
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()
	orderCollection.InsertOne(context.Background(), order)
	return order.OrderId
}
