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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

var orderItemCollection = database.OpenConnection("orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		result, err := orderItemCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		var allOrderItems []bson.M
		err = result.All(ctx, &allOrderItems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}

		c.JSON(200, gin.H{
			"status":  200,
			"message": "Success",
			"data":    allOrderItems,
		})
	}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()

		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"orderitemid": orderItemId}).Decode(&orderItem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"STatus":  http.StatusInternalServerError,
				"error":   err,
				"message": "Internal Server Error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": orderItem,
		})
	}
}
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occur while listing order items by order id",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": allOrderItems,
		})
	}
}
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx = context.Background()

	matchStage := primitive.D{
		primitive.E{Key: "$match", Value: bson.M{"orderid": id}},
	}
	lookupFoodStage := primitive.D{
		primitive.E{Key: "$lookup", Value: bson.M{
			"from":         "food",
			"localField":   "foodid",
			"foreignField": "foodid",
			"as":           "food",
		}},
	}

	unwindFoodStage := primitive.D{
		primitive.E{Key: "$unwind", Value: bson.M{
			"path":                       "$food",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	lookupOrderStage := primitive.D{
		primitive.E{Key: "$lookup", Value: bson.M{
			"from":         "order",
			"localField":   "orderid",
			"foreignField": "orderid",
			"as":           "order",
		}},
	}
	unwindOrderStage := primitive.D{
		primitive.E{Key: "$unwind", Value: bson.M{
			"path":                       "$order",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	lookupTableStage := primitive.D{
		primitive.E{Key: "$lookup", Value: bson.M{
			"from":         "table",
			"localField":   "tableid",
			"foreignField": "tableid",
			"as":           "table",
		}},
	}
	unwindTableStage := primitive.D{
		primitive.E{Key: "$unwind", Value: bson.M{
			"path":                       "$table",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	projectStage := bson.D{
		primitive.E{Key: "$project", Value: bson.M{
			"id":          1,
			"amount":      "$food.price",
			"totalcount":  1,
			"foodname":    "$food.name",
			"foodimage":   "$food.foodimage",
			"tablenumber": "$table.tablenumber",
			"tableid":     "$table.tableid",
			"orderid":     "$order.orderid",
			"price":       "$food.price",
			"quantity":    1,
		}},
	}
	groupStage := primitive.D{
		primitive.E{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"orderid":     "$orderid",
				"tableid":     "$tableid",
				"tablenumber": "$tablenumber",
			},
			"paymentdue": bson.M{
				"$sum": "$amount",
			},
			"totalcount": bson.M{
				"$sum": 1,
			},
			"orderitems": bson.M{
				"$push": "$$ROOT",
			},
		}},
	}
	projectStage2 := bson.D{
		primitive.E{Key: "$project", Value: bson.M{
			"id":          1,
			"paymentdue":  1,
			"totalcount":  1,
			"orderitems":  1,
			"tablenumber": "$_id.tablenumber",
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupFoodStage,
		unwindFoodStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	return OrderItems, err

}
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var orderItemPack OrderItemPack
		var order models.Order
		err := c.BindJSON(&orderItemPack)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})

			return
		}
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeInserted := []interface{}{}
		order.TableId = orderItemPack.TableId
		orderId := orderItemOrderCreator(order)
		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderId = orderId
			validateErr := validate.Struct(orderItem)
			if validateErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": validateErr.Error(),
				})

				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemId = orderItem.ID.Hex()
			num := toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}
		insertedOrderItem, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": insertedOrderItem,
		})
	}
}
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D
		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"order_item_id": orderItem.OrderItemId}})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"quantity": orderItem.Quantity}})
		}
		if orderItem.FoodId != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"food_id": orderItem.FoodId}})
		}
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": orderItem.UpdatedAt}})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		_, err := orderItemCollection.UpdateOne(ctx, filter, updateObj, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			return
		}

	}
}
