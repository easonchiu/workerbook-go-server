package model

import (
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const GroupCollection = "groups"

// collection schema
type Group struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 分组名
  Name string `json:"name"`

  // 分组下的用户数
  Count int `json:"count"`

  // 创建时间
  CreateTime time.Time `json:"createTime" bson:"createTime"`
}
