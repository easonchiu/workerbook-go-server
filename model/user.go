package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// collection name
const UserCollection = "users"

// collection schema
type User struct {
	// id
	Id bson.ObjectId `json:"id" bson:"_id"`

	// 昵称
	NickName string `json:"nickname"`

	// 邮箱
	Email string `json:"email"`

	// 用户名
	UserName string `json:"username"`

	// 分组id
	Gid string `json:"gid"`

  // 分组名
  GroupName string `json:"groupName" bson:"groupName"`

	// 手机号
	Mobile string `json:"mobile"`

	// 密码
	Password string `json:"password"`

	// 1: 普通用户， 2: 部门Leader，99: 管理员
	Role int `json:"role"`

	// 创建时间
	CreateTime time.Time `json:"createTime" bson:"createTime"`

	// 用户是否存在
	Exist bool `json:"exist"`
}

// user result
type UserResult struct {
  // id
	Id bson.ObjectId `json:"id" bson:"_id"`

	// 昵称
	NickName string `json:"nickname"`

	// 分组id
	Gid string `json:"gid"`

  // 分组名
  GroupName string `json:"groupName" bson:"groupName"`

  // 邮箱
  Email string `json:"email"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
	Role int `json:"role"`

	// 创建时间
	CreateTime time.Time `json:"createTime" bson:"createTime"`
}
