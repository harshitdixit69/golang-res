package controller

import (
	"context"
	"fmt"
	"log"
	"restaurantProject/database"
	"restaurantProject/helpers"
	"restaurantProject/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenConnection("user")
var ctx = context.Background()

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := userCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(200, gin.H{
				"message": err,
			})
			return
		}
		var resultData []bson.M
		err = result.All(ctx, &resultData)
		if err != nil {
			c.JSON(200, gin.H{
				"message": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    resultData,
		})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		err := userCollection.FindOne(ctx, primitive.E{Key: "$set", Value: bson.M{"user_id": userId}}).Decode(&user)
		if err != nil {
			c.JSON(404, gin.H{
				"message": "User not found",
			})
			log.Println(err)

			return
		}
		c.JSON(200, gin.H{
			"message": "User found",
			"user":    user,
		})
	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err,
			})
			log.Println(err)
			return
		}
		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(400, gin.H{
				"message": err,
			})
			log.Println(err)
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Something went wrong while getting user count",
			})
			log.Println(err)
			return
		}
		if count > 0 {
			c.JSON(400, gin.H{
				"message": "Email already exists",
			})
			log.Println(err)
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Something went wrong while getting user count",
			})
			log.Println(err)
			return
		}
		if count > 0 {
			c.JSON(400, gin.H{
				"message": "Phone No already exists",
			})
			log.Println(err)
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, *user.LastName, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Something went wrong while creating user",
			})
			fmt.Println("6")
			log.Println(err)
			return
		}
		c.JSON(200, gin.H{
			"message": "User created",
			"data":    result,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err,
			})

			return
		}
		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "User not found",
			})

			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(400, gin.H{
				"message": msg,
			})

			return
		}
		token, refreshToken, err := helpers.GenerateAllToken(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Something went wrong",
			})
			return
		}
		c.SetCookie("tokens", token, 10, "/", "localhost", false, true)

		helpers.UpdateAllToken(token, refreshToken, foundUser.UserId)
		c.JSON(200, gin.H{
			"message": "User found",
			"data":    foundUser,
		})
	}
}
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	if err != nil {
		return false, "Wrong password"
	}
	return true, ""
}
