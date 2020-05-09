package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	hosts              = "localhost:27017"
	database           = "vpackagerdb"
	username           = ""
	password           = ""
	contentsCollection = "contents"
	streamsCollection  = "streams"
)

type Mongodb struct {
	session *mgo.Session
}

func NewMongodb() *Mongodb {
	return &Mongodb{
		session: nil,
	}
}

func (db *Mongodb) Init() error {

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return err
	}
	db.session = session

	return nil
}

func (db *Mongodb) SaveContent(content Content) error {
	if db.session == nil {
		return fmt.Errorf("mongodb not initialized")
	}

	collection := db.session.DB(database).C(contentsCollection)

	err := collection.Insert(content)
	if err != nil {
		return err
	}
	return nil
}

func (db *Mongodb) GetContentID(content Content) (string, error) {
	if db.session == nil {
		return "", fmt.Errorf("mongodb not initialized")
	}

	collection := db.session.DB(database).C(contentsCollection)

	var newContent Content
	err := collection.Find(bson.M{"name": content.Name}).One(&newContent)
	if err != nil {
		return "", err
	}
	return newContent.ID.Hex(), nil
}
