package model

import (
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const ProjectCollection = "projects"

// collection schema
type Project struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 项目名
  Name string `json:"name"`

  // 状态 1. 启用 2. 已归档
  Status int `json:"status"`

  // 创建时间
  CreateTime time.Time `json:"createTime" bson:"createTime"`
}