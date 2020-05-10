package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Content struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name"`
}

type Stream struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ContentID string             `json:"content_id"`
	Key       string             `json:"key"`
	Kid       string             `json:"kid"`
	Status    string             `json:"status"`
	Url       string             `json:"url"`
}
