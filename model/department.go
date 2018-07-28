package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/util"
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

  // 操作人
  Editor mgo.DBRef `bson:"editor,omitempty"`

  // 操作时间
  EditTime time.Time `bson:"editTime,omitempty"`
}

func (d Department) GetMap(forgets ... string) gin.H {
  data := gin.H{
    "id":         d.Id,
    "name":       d.Name,
    "userCount":  d.UserCount,
    "createTime": d.CreateTime,
  }

  util.Forget(data, forgets...)

  return data
}

// 部门列表结构
type DepartmentList struct {
  List  *[]Department
  Count int
  Limit int
  Skip  int
}

// 列表的迭代器
func (d DepartmentList) Each(fn func(Department) gin.H) gin.H {
  data := gin.H{}

  if d.Limit != 0 {
    data = gin.H{
      "count": d.Count,
      "limit": d.Limit,
      "skip":  d.Skip,
    }
  }

  var list []gin.H
  for _, item := range *d.List {
    list = append(list, fn(item))
  }

  data["list"] = list

  return data
}
