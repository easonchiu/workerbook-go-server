package model

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// collection name
const UserAnalyticsCollection = "user_analytics"

type UserAnalyticsDaily struct {
  MissionId   bson.ObjectId `bson:"missionId"`
  MissionName string        `bson:"missionName"`
  Progress    int           `bson:"progress"`
  IsTimeout   bool          `bson:"isTimeout"`
}

// collection schema
type UserAnalytics struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // department
  User mgo.DBRef `bson:"user"`

  // department
  Department mgo.DBRef `bson:"department"`

  // 日报数据
  Dailies []UserAnalyticsDaily `bson:"daily"`

  // 日期(20180201这样的格式)
  Day string `bson:"day"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`
}
