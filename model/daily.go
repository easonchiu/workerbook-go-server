package model

import (
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const DailyCollection = "dailies"

type DailyMission struct {
  // 项目归属id
  ProjectId string `json:"projectId"`

  // 任务归属id
  MissionId string `json:"missionId"`

  // 项目名
  ProjectName string `json:"projectName"`

  // 任务名
  MissionName string `json:"missionName"`

  // 任务类型 1. 评估 2. 开发 3. 测试 4. 上线
  MissionType int `json:"missionType"`

  // 任务类型中文字
  MissionTypeText string `json:"missionTypeText"`

  // 我对该任务的进度，这里的进度和任务进度不是同一个进度
  MyProgress int `json:"myProgress"`
}

// daily schema
type DailyItem struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 内容
  Content string `json:"record" bson:"record"`

  // 进度
  Progress int `json:"progress"`

  // 任务(如果是任务日报，进度使用的是任务的进度)
  Mission DailyMission `json:"mission"`
}

// collection schema
type Daily struct {
  // id
  Id bson.ObjectId `json:"id" bson:"_id"`

  // 用户id
  Uid string `json:"uid"`

  // 用户的分组id
  GroupId string `json:"groupId"`

  // 日期
  Day string `json:"day"`

  // 日报数据
  DailyList []DailyItem `json:"dailyList" bson:"dailyList"`

  // 发布时间
  CreateTime time.Time `json:"createTime"`

  // 更新时间
  UpdateTime time.Time `json:"updateTime"`
}


