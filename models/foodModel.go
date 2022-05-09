package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      *string            `json:"name" validate:"required,min=3,max=50"`
	Price     *float64           `json:"price" validate:"required,min=0"`
	FoodImage *string            `json:"food_image" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	FoodId    string             `json:"food_id"`
	MenuId    *string            `json:"menu_id" validate:"required"`
}
