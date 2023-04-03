package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	ifIDExistsDB(string) error
	insertDB(*User) error
	getDB(string, *User) error
	updateIdDB(string, string) error
	updatePasswordDB(string, string) error
	updateCanvasDB(string, string) error
	deleteDB(string) error
}

type MongoStore struct {
	collection *mongo.Collection
}

func newDBStore() (*MongoStore, error) {
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

func (m *MongoStore) ifIDExistsDB(s string) error {
	result := m.collection.FindOne(context.TODO(), bson.M{"_id": s})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (m *MongoStore) insertDB(u *User) error {
	value := bson.D{{Key: "_id", Value: u.ID}, {Key: "password", Value: u.Password}, {Key: "canvas_data", Value: u.CanvasData}}
	_, err := m.collection.InsertOne(context.TODO(), value)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoStore) getDB(s string, resp *User) error {
	err := m.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: s}}).Decode(&resp)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoStore) updatePasswordDB(id string, pass string) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: pass}}}}
	_, err := m.collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, update)

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStore) updateCanvasDB(id string, c string) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "canvas_data", Value: c}}}}
	_, err := m.collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, update)

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStore) deleteDB(s string) error {
	_, err := m.collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: s}})
	if err != nil {
		return err
	}
	return nil
}
