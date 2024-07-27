package models

type User struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Gender string `json:"gender" bson:"gender"`
	Age    int    `json:"age" bson:"age"`
}
