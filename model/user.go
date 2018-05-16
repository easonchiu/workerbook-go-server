package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// collection name
const UserCollection = "users"

// collection schema
type User struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	// 昵称
	NickName string `json:"nickname"`
	// 邮箱
	Email string `json:"email"`
	// 用户名
	UserName string `json:"username"`
	// 分组id
	Gid string `json:"gid"`
	// 手机号
	Mobile string `json:"mobile"`
	// 密码
	Password string `json:"password"`
	// 1：管理员， 2：普通用户
	Role int `json:"role"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
}

// user result
type UserResult struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	// 昵称
	NickName string `json:"nickname"`
	// 分组id
	Gid string `json:"gid"`
	// 1：管理员， 2：普通用户
	Role int `json:"role"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
}
