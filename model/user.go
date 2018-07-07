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
  NickName string `json:"nickname"`

  // 邮箱
  Email string `json:"email"`

  // 用户名
  UserName string `json:"username"`

  // 部门id
  DepartmentId string `json:"departmentId"`

  // 手机号
  Mobile string `json:"mobile"`

  // 密码
  Password string `json:"password"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
  Role int `json:"role"`

  // 创建时间
  CreateTime time.Time `json:"createTime"`

  // 用户状态
  Status int `json:"status"`
}

// user result
type UserResult struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 昵称
  NickName string `json:"nickname"`

  // 部门id
  DepartmentId string `json:"departmentId"`

  // 邮箱
  Email string `json:"email"`

  // 1: 普通用户， 2: 部门Leader，99: 管理员
  Role int `json:"role"`
}
