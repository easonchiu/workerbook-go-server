package model

import "gopkg.in/mgo.v2/bson"

const UserCollection = "users"

type User struct {
	A		string			`json:"a" bson:"a"`
	Id		bson.ObjectId	`json:"id,omitempty" bson:"_id,omitempty"`
}