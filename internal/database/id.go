package database

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewStringID() string {
	return primitive.NewObjectID().Hex()
}
