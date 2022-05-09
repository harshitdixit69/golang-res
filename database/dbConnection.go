package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	MongoDb := "mongodb+srv://harshit65:harshit65@cluster0.vggap.mongodb.net/"
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}
	var ctx = context.Background()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("1")
		log.Fatal(err)
	}
	return client
}

var Client *mongo.Client = DBInstance()

func OpenConnection(collectionName string) *mongo.Collection {
	var collection *mongo.Collection = Client.Database("restaurants").Collection(collectionName)
	return collection
}
