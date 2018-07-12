package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  `time`
)

// collection name
const DepartmentCollection = "departments"

// collection schema
type Department struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 部门名
  Name string `bson:"name,omitempty"`

  // 部门下的用户数
  UserCount int `bson:"userCount,omitempty"`

  // 创建时间
  CreateTime time.Time `bson:"createTime,omitempty"`
}

func (d Department) GetMap(db *mgo.Database) gin.H {
  return gin.H{
    "id":         d.Id,
    "name":       d.Name,
    "userCount":  d.UserCount,
    "createTime": d.CreateTime,
  }
}
