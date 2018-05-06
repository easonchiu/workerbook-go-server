package model

import "gopkg.in/mgo.v2/bson"

// collection name
const UserCollection = "users"

// collection model
type User struct {
	A		string			`json:"a" bson:"a"`
	Id		bson.ObjectId	`json:"id,omitempty" bson:"_id,omitempty"`
}