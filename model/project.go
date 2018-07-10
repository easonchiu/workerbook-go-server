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
  Name string `bson:"name"`

  // 状态 1. 启用 2. 已归档 3. 停止
  Status int `bson:"status"`

  // 截至时间
  Deadline time.Time `bson:"deadline"`

  // 参与的部门
  Departments []mgo.DBRef `bson:"departments"`

  // 任务列表
  Missions []Mission `bson:"missions"`

  // 项目说明
  Description string `bson:"description"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`

  // 进度
  Progress int `bson:"progress"`
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
  return gin.H{
    "id": p.Id,
    "name": p.Name,
    "missionCount": len(p.Missions),
    "deadline": p.Deadline,
    "description": p.Description,
    "createTime": p.CreateTime,
    "progress": p.Progress,
    "departments": departments,
  }
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
