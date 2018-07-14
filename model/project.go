package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const ProjectCollection = "projects"

// collection schema
type Project struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 项目名
  Name string `bson:"name,omitempty"`

  // 状态 1. 启用 2. 已归档 3. 停止
  Status int `bson:"status,omitempty"`

  // 截至时间
  Deadline time.Time `bson:"deadline,omitempty"`

  // 参与的部门
  Departments []mgo.DBRef `bson:"departments,omitempty"`

  // 项目说明
  Description string `bson:"description,omitempty"`

  // 创建时间
  CreateTime time.Time `bson:"createTime,omitempty"`

  // 进度
  Progress int `bson:"progress,omitempty"`

  // 权重 1. 红(紧急) 2. 黄(重要) 3. 绿(一般)
  Weight int `bson:"weight,omitempty"`

  // 任务
  Missions []mgo.DBRef `bson:"missions,omitempty"`

  // 是否存在
  Exist bool `bson:"exist"`
}

func (p Project) GetMap(db *mgo.Database) gin.H {
  var departments []gin.H
  for _, item := range p.Departments {
    department := new(Department)
    err := db.FindRef(&item).One(department)
    if err == nil {
      departments = append(departments, department.GetMap(db))
    }
  }
  var missions []gin.H
  for _, item := range p.Missions {
    mission := new(Mission)
    err := db.FindRef(&item).One(mission)
    if err == nil {
      missions = append(missions, mission.GetMap(db))
    }
  }
  return gin.H{
    "id":          p.Id,
    "name":        p.Name,
    "deadline":    p.Deadline,
    "description": p.Description,
    "createTime":  p.CreateTime,
    "progress":    p.Progress,
    "departments": departments,
    "missions":    missions,
    "weight":      p.Weight,
  }
}
