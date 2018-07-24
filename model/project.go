package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/util"
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

func (p Project) GetMap(forgets ... string) gin.H {
  data := gin.H{
    "id":          p.Id,
    "name":        p.Name,
    "deadline":    p.Deadline,
    "description": p.Description,
    "createTime":  p.CreateTime,
    "progress":    p.Progress,
    "weight":      p.Weight,
  }

  // departments
  var departments []gin.H

  for _, item := range p.Departments {
    if item.Id != "" {
      departments = append(departments, gin.H{
        "id": item.Id.(bson.ObjectId),
      })
    }
  }

  data["departments"] = departments

  // missions
  var missions []gin.H

  for _, item := range p.Missions {
    if item.Id != "" {
      missions = append(missions, gin.H{
        "id": item.Id.(bson.ObjectId),
      })
    }
  }

  data["missions"] = missions

  util.Forget(data, forgets...)

  return data
}

// 项目列表结构
type ProjectList struct {
  List  *[]Project
  Count int
  Limit int
  Skip  int
}

// 列表的迭代器
func (d ProjectList) Each(fn func(Project) gin.H) gin.H {
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
