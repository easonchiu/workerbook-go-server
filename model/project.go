package model

import (
  "github.com/gin-gonic/gin"
  "github.com/influxdata/influxdb/pkg/slices"
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
  Name string `bson:"name"`

  // 状态 1. 启用 2. 已归档 3. 停止
  Status int `bson:"status"`

  // 截至时间
  Deadline time.Time `bson:"deadline"`

  // 参与的部门
  Departments []mgo.DBRef `bson:"departments"`

  // 项目说明
  Description string `bson:"description"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`

  // 进度
  Progress int `bson:"progress"`

  // 权重 1. 红(紧急) 2. 黄(重要) 3. 绿(一般)
  Weight int `bson:"weight"`

  // 任务
  Missions []mgo.DBRef `bson:"missions"`

  // 是否存在
  Exist bool `bson:"exist"`
}

func (p Project) GetMap(db *mgo.Database, refs ... string) gin.H {
  data := gin.H{
    "id":          p.Id,
    "name":        p.Name,
    "deadline":    p.Deadline,
    "description": p.Description,
    "createTime":  p.CreateTime,
    "progress":    p.Progress,
    "weight":      p.Weight,
  }

  // departments refs
  var departments []gin.H

  if slices.Exists(refs, "departments") {
    for _, item := range p.Departments {
      department := new(Department)
      err := db.FindRef(&item).One(department)
      if err == nil {
        departments = append(departments, department.GetMap(db))
      }
    }
  } else {
    for _, item := range p.Departments {
      if item.Id != "" {
        departments = append(departments, gin.H{
          "id": item.Id.(bson.ObjectId),
        })
      }
    }
  }

  data["departments"] = departments

  // missions refs
  var missions []gin.H

  if slices.Exists(refs, "missions") {
    for _, item := range p.Missions {
      mission := new(Mission)
      err := db.FindRef(&item).One(mission)
      if err == nil {
        missions = append(missions, mission.GetMap(db, "user"))
      }
    }
  } else {
    for _, item := range p.Missions {
      if item.Id != "" {
        missions = append(missions, gin.H{
          "id": item.Id.(bson.ObjectId),
        })
      }
    }
  }

  data["missions"] = missions

  return data
}
