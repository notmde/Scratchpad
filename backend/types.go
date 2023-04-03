package main

type User struct {
	ID         string `json:"_id" bson:"_id"`
	Password   string `json:"password" bson:"password"`
	CanvasData string `json:"canvas_data" bson:"canvas_data"`
}
