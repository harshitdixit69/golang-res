package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID            primitive.ObjectID `bson:"_id"`
	NumberOfGuest *int64             `json:"number_of_guest" validate:"required,min=1,max=10"`
	TableNumber   *int64             `json:"table_number" validate:"required,min=1,max=10"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	TableId       string             `json:"table_id"`
}
