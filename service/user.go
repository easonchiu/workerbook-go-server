package service

import (
	"github.com/gin-gonic/gin"
	"workerbook/db"
	"workerbook/model"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

// Insert user info into database.
func CreateUser(data model.User) error {
	db, close, err := db.CloneDB()

	if err != nil {
		return err
	} else {
		defer close()
	}

	// check the data is error or not.
	if data.UserName == "" {
		return errors.New("用户名不能为空")
	} else if data.Password == "" {
		return errors.New("密码不能为空")
	} else if data.NickName == "" {
		return errors.New("昵称能为空")
	} else if !bson.IsObjectIdHex(data.Gid) {
		return errors.New("分组号错误")
	}

	// username must be the only.
	count, err := db.C(model.UserCollection).Find(bson.M{"username": data.UserName}).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("已存在相同的用户")
	}

	// create a new object id.
	if data.Id == "" {
		data.Id = bson.NewObjectId()
	}

	// insert it.
	err = db.C(model.UserCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// Query user info by id.
func GetUserInfoById(id bson.ObjectId) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := model.User{}

	err = db.C(model.UserCollection).FindId(id).One(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"data": data,
	}, nil
}

// Query users list with skip and limit.
func GetUsersList(skip int, limit int) (gin.H, error) {
	db, close, err := db.CloneDB()

	if err != nil {
		return nil, err
	} else {
		defer close()
	}

	data := make([]model.User, limit)

	err = db.C(model.UserCollection).Find(bson.M{}).Skip(skip).Limit(limit).All(&data)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list": data,
	}, nil
}