package model

import (
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const DepartmentCollection = "departments"

// collection schema
type Department struct {
  // 部门名
  Name string `json:"name"`

  // 部门下的用户数
  UserCount int `json:"userCount"`

  // 创建时间
  CreateTime time.Time `json:"createTime"`
}


type DepartmentResult struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 部门名
  Name string `json:"name"`

  // 部门下的用户数
  UserCount int `json:"userCount"`
}
