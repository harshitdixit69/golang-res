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

type InvoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	OrderDetails   interface{}
	PaymentDueDate time.Time
}

var invoiceCollection = database.OpenConnection("invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		result, err := invoiceCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		var allInvoices []bson.M
		err = result.All(ctx, &allInvoices)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    allInvoices,
		})
	}
}
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invoice not found",
			})
		}
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invoice not found",
			})
		}
		invoiceView.OrderId = invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}
		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"]
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success",
			"data":    invoiceView,
		})
	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		var invoice models.Invoice

		bindJsonErr := c.BindJSON(&invoice)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})

			return
		}
		validateErr := validate.Struct(invoice)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})

			return
		}
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "order not found",
			})

			return
		}
		status := "PENDING"
		if invoice.PaymentStatus != nil {
			invoice.PaymentStatus = &status
		}
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()
		insertResult, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
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
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = context.Background()
		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice
		bindJsonErr := c.BindJSON(&invoice)
		if bindJsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": bindJsonErr.Error(),
			})

			return
		}
		filter := bson.M{"invoice_id": invoiceId}
		var updateObj primitive.D
		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"payment_method": *invoice.PaymentMethod}})
		}
		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"payment_status": *invoice.PaymentStatus}})
		}

		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, primitive.E{Key: "$set", Value: bson.M{"updated_at": invoice.UpdatedAt}})
		upsertL := true
		opt := options.UpdateOptions{
			Upsert: &upsertL,
		}
		status := "PENDING"
		if invoice.PaymentStatus != nil {
			invoice.PaymentStatus = &status
		}
		result, err := invoiceCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj}, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
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
func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
