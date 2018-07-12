package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  `gopkg.in/mgo.v2/bson`
  "time"
)

// collection name
const MissionCollection = "missions"

// collection schema
type Mission struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 任务名
  Name string `bson:"name,omitempty"`

  // 进度，这个值的算法为：于次日10点结算前一天的所有日报，该天如果有用户针对该任务写日报，则每人针该进度写的值取平均
  Progress int `bson:"progress,omitempty"`

  // 描述
  Description string `bson:"description,omitempty"`

  // 截至时间
  Deadline time.Time `bson:"deadline,omitempty"`

  // 状态 1. 正常 2. 停止
  Status int `bson:"status,omitempty"`

  // 创建时间
  CreateTime time.Time `bson:"createTime,omitempty"`
}

func (m Mission) GetMap(db *mgo.Database) gin.H {
  return gin.H{
    "id":          m.Id,
    "name":        m.Name,
    "description": m.Description,
    "createTime":  m.CreateTime,
    "deadline":    m.Deadline,
    "status":      m.Status,
    "progress":    m.Progress,
  }
}
