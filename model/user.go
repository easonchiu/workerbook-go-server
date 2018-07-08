package model

import (
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const UserCollection = "users"

// collection schema
type User struct {
  // 昵称
  NickName string `json:"nickname" bson:"nickname"`

  // 邮箱
  Email string `json:"email" bson:"email"`

  // 用户名
  UserName string `json:"username" bson:"username"`

  // 部门id
  DepartmentId string `json:"departmentId" bson:"departmentId"`

  // 手机号
  Mobile string `json:"mobile" bson:"mobile"`

  // 密码
  Password string `json:"password" bson:"password"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
  Role int `json:"role" bson:"role"`

  // 创建时间
  CreateTime time.Time `json:"createTime" bson:"createTime"`

  // 用户状态
  Status int `json:"status" bson:"status"`
}

// user result
type UserResult struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 昵称
  NickName string `json:"nickname" bson:"nickname"`

  // 部门id
  DepartmentId string `json:"departmentId" bson:"departmentId"`

  // 邮箱
  Email string `json:"email" bson:"email"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
  Role int `json:"role" bson:"role"`
}

// user result at console
type UserConsoleResult struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 昵称
  NickName string `json:"nickname" bson:"nickname"`

  // 邮箱
  Email string `json:"email" bson:"email"`

  // 用户名
  UserName string `json:"username" bson:"username"`

  // 部门id
  DepartmentId string `json:"departmentId" bson:"departmentId"`

  // 手机号
  Mobile string `json:"mobile" bson:"mobile"`

  // 密码
  Password string `json:"password" bson:"password"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
  Role int `json:"role" bson:"role"`

  // 创建时间
  CreateTime time.Time `json:"createTime" bson:"createTime"`

  // 用户状态
  Status int `json:"status" bson:"status"`
}