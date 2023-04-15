package storage

import (
	"context"
	"os"

	"github.com/scratchpad-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	IfIDExistsDB(string) error
	InsertDB(*types.User) error
	GetDB(string, *types.User) error
	UpdateCanvasDB(string, string) error
}

type MongoStore struct {
	collection *mongo.Collection
}

func NewDBStore() (*MongoStore, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return &MongoStore{
		collection: client.Database("Scratchpad").Collection("canvasData"),
	}, nil
}

func (m *MongoStore) IfIDExistsDB(s string) error {
	result := m.collection.FindOne(context.TODO(), bson.M{"_id": s})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (m *MongoStore) InsertDB(u *types.User) error {
	value := bson.D{{Key: "_id", Value: u.ID}, {Key: "password", Value: u.Password}, {Key: "canvas_data", Value: u.CanvasData}}
	_, err := m.collection.InsertOne(context.TODO(), value)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoStore) GetDB(s string, resp *types.User) error {
	err := m.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: s}}).Decode(&resp)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoStore) UpdateCanvasDB(id string, c string) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "canvas_data", Value: c}}}}
	_, err := m.collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, update)

	if err != nil {
		return err
	}

	return nil
}
