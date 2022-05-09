package helpers

import (
	"context"
	"os"
	"restaurantProject/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDecodes struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenConnection("user")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllToken(email string, firstName string, lastName string, uid string) (string, string, error) {
	claim := &SignedDecodes{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}
	refreshClaims := &SignedDecodes{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func UpdateAllToken(signedToken string, signedrefreshToken string, userId string) {
	var ctx = context.Background()
	var updateObj primitive.D
	updateObj = append(updateObj, primitive.E{Key: "$set", Value: primitive.E{
		Key: "token", Value: signedToken,
	}})
	updateObj = append(updateObj, primitive.E{Key: "$set", Value: primitive.E{
		Key: "refreshToken", Value: signedrefreshToken,
	}})
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "$set", Value: primitive.E{Key: "updatedAt", Value: updateAt}})
	upsert := true
	filter := bson.D{primitive.E{Key: "user_id", Value: userId}}
	opt := options.UpdateOptions{Upsert: &upsert}
	_, err := userCollection.UpdateOne(ctx, filter, updateObj, &opt)

	if err != nil {
		return
	}

}

func ValidateToken(signedToken string) (SignedDecodes, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDecodes{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return SignedDecodes{}, "Cannot validate token"
	}
	claim, ok := token.Claims.(*SignedDecodes)
	if !ok || !token.Valid {
		return SignedDecodes{}, "Token is invalid"
	}
	if claim.ExpiresAt < time.Now().Local().Unix() {
		return SignedDecodes{}, "token is expired"
	}
	return *claim, ""
}
