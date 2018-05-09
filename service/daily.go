package service

import (
	"github.com/gin-gonic/gin"
	"workerbook/db"
	"workerbook/model"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Insert daily info into database.
func CreateDaily(data model.Daily) error {
	db, close, err := db.CloneDB()

	if err != nil {
		return err
	} else {
		defer close()
	}

	// check the data is error or not.


	// supplement other data.
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()

	// create a new object id.
	if data.Id == "" {
		data.Id = bson.NewObjectId()
	}

	// insert it.
	err = db.C(model.DailyCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// Query daily info by id.
func GetDailyInfoById(id bson.ObjectId) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := model.Daily{}

	err = db.C(model.DailyCollection).FindId(id).One(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"data": data,
	}, nil
}

// Query dailies list with skip and limit.
func GetDailiesList(skip int, limit int) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := make([]model.Daily, limit)

	err = db.C(model.DailyCollection).Find(bson.M{}).Skip(skip).Limit(limit).All(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list": data,
	}, nil
}