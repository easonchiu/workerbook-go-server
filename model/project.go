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

  // 状态 1. 启用 2. 已归档 3. 停止
  Status int `json:"status"`

  // 截至时间
  Expire time.Time `json:"expire"`

  // 任务列表
  MissionList []Mission `json:"missionList"`

  // 项目说明
  Summary string `json:"summary"`

  // 创建时间
  CreateTime time.Time `json:"createTime"`
}


// collection schema
type Mission struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 任务名
  Name string `json:"name"`

  // 进度，这个值的算法为：于次日10点结算前一天的所有日报，该天如果有用户针对该任务写日报，则每人针该进度写的值取平均
  Progress int `json:"progress"`

  // 描述
  Summary string `json:"summary"`

  // 类型 1. 评估 2. 开发 3. 测试 4. 上线
  Type int `json:"type"`

  // 类型中文字
  TypeText string `json:"typeText"`

  // 状态 1. 正常 2. 停止
  Status int `json:"status"`
}