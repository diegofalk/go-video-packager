package database

import "gopkg.in/mgo.v2/bson"

type Content struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name 	string `json:"name"`
}

type Stream struct {
	ContentID string `json:"content_id"`
	Key       string `json:"key"`
	Kid       string `json:"kid"`
	Status    string `json:"status"`
}
