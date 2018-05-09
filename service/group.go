package service

import (
	"github.com/gin-gonic/gin"
	"workerbook/db"
	"workerbook/model"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"time"
)

// Insert group info into database.
func CreateGroup(data model.Group) error {
	db, close, err := db.CloneDB()

	if err != nil {
		return err
	} else {
		defer close()
	}

	// check the data is error or not.
	if data.Name == "" {
		return errors.New("分组不能为空")
	}

	// supplement other data.
	data.CreateTime = time.Now()
	data.Count = 0

	// name must be the only.
	count, err := db.C(model.GroupCollection).Find(bson.M{"name": data.Name}).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("已存在相同的分组")
	}

	// create a new object id.
	if data.Id == "" {
		data.Id = bson.NewObjectId()
	}

	// insert it.
	err = db.C(model.GroupCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// Query group info by id.
func GetGroupInfoById(id bson.ObjectId) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := model.Group{}

	err = db.C(model.GroupCollection).FindId(id).One(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"data": data,
	}, nil
}

// Query groups list with skip and limit.
func GetGroupsList(skip int, limit int) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := make([]model.Group, limit)

	err = db.C(model.GroupCollection).Find(bson.M{}).Skip(skip).Limit(limit).All(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list": data,
	}, nil
}