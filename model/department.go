package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// collection name
const DepartmentCollection = "departments"

// collection schema
type Department struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 部门名
  Name string `bson:"name"`

  // 部门下的用户数
  UserCount int `bson:"userCount"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`

  // 是否存在
  Exist bool `bson:"exist"`
}

func (d Department) GetMap(refFunc func(mgo.DBRef) (gin.H, bool), refs ... string) gin.H {
  data := gin.H{
    "id":         d.Id,
    "name":       d.Name,
    "userCount":  d.UserCount,
    "createTime": d.CreateTime,
  }
  return data
}
