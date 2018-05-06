package service

import (
	"github.com/gin-gonic/gin"
	"web/db"
	"web/model"
	"gopkg.in/mgo.v2/bson"
)

func GetUserInfoById(id string) (gin.H, error) {

	mongo := db.MuseDB()

	data := []model.User{}
	err := mongo.C(model.UserCollection).Find(bson.M{}).All(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list": data,
	}, nil
}