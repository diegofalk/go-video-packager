package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri                = "mongodb://localhost:27017"
	database           = "vpackagerdb"
	contentsCollection = "contents"
	streamsCollection  = "streams"
)

type Mongodb struct {
	client *mongo.Client
}

func NewMongodb() *Mongodb {
	return &Mongodb{
		client: nil,
	}
}

func (db *Mongodb) Init() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return err
	}
	db.client = client

	return nil
}

func (db *Mongodb) InsertContent(content Content) (string, error) {
	if db.client == nil {
		return "", fmt.Errorf("mongodb not initialized")
	}

	collection := db.client.Database(database).Collection(contentsCollection)

	result, err := collection.InsertOne(context.Background(), content)
	if err != nil {
		return "", err
	}
	contentId := result.InsertedID.(primitive.ObjectID).Hex()

	return contentId, nil
}

func (db *Mongodb) GetContent(id string) (Content, error) {
	if db.client == nil {
		return Content{}, fmt.Errorf("mongodb not initialized")
	}

	collection := db.client.Database(database).Collection(contentsCollection)

	var content Content
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Content{}, err
	}
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(context.Background(), filter).Decode(&content)
	if err != nil {
		return Content{}, err
	}
	return content, nil
}

func (db *Mongodb) InsertStream(stream Stream) (string, error) {
	if db.client == nil {
		return "", fmt.Errorf("mongodb not initialized")
	}

	collection := db.client.Database(database).Collection(streamsCollection)

	result, err := collection.InsertOne(context.Background(), stream)
	if err != nil {
		return "", err
	}
	streamId := result.InsertedID.(primitive.ObjectID).Hex()

	return streamId, nil
}