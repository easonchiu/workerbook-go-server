package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// collection name
const DailyCollection = "dailies"

// daily schema
type DailyItem struct {
  // id
  Id bson.ObjectId `bson:"_id"`

  // 内容
  Content string `bson:"content"`

  // 进度
  Progress int `bson:"progress"`

  // 任务名称
  MissionName string `bson:"missionName"`

  // 任务id
  MissionId bson.ObjectId `bson:"missionId"`

  // 项目名称
  ProjectName string `bson:"projectName"`

  // 项目id
  ProjectId bson.ObjectId `bson:"projectId"`
}

// collection schema
type Daily struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 用户信息
  User mgo.DBRef `bson:"user"`

  // 部门名称(因为部门有可能删除，且这里只需要名称，所以直接存下来)
  DepartmentName string `bson:"departmentName"`

  // 日期(20180201这样的格式)
  Day string `bson:"day"`

  // 日报数据
  Dailies []*DailyItem `bson:"dailies"`

  // 发布时间
  CreateTime time.Time `bson:"createTime"`

  // 更新时间
  UpdateTime time.Time `bson:"updateTime"`
}

func (d *Daily) GetMap(forgets ... string) gin.H {
  data := gin.H{
    "id": d.Id,
    "user": gin.H{
      "id": d.User.Id,
    },
    "departmentName": d.DepartmentName,
    "createTime":     d.CreateTime,
    "updateTime":     d.UpdateTime,
  }

  var dailies []gin.H
  for _, item := range d.Dailies {
    dailies = append(dailies, gin.H{
      "id":       item.Id,
      "content":  item.Content,
      "progress": item.Progress,
      "mission": gin.H{
        "name": item.MissionName,
        "id":   item.MissionId,
      },
      "project": gin.H{
        "name": item.ProjectName,
        "id":   item.ProjectId,
      },
    })
  }

  data["dailies"] = dailies

  if forgets != nil {
    if forgets[0] == REMEMBER {
      remember(data, forgets[1:]...)
    } else {
      forget(data, forgets...)
    }
  }

  return data
}

// 日报列表结构
type DailyList struct {
  List  []*Daily
  Count int
  Limit int
  Skip  int
}

// 列表的迭代器
func (d *DailyList) Each(fn func(*Daily) gin.H) gin.H {
  data := gin.H{}

  if d.Limit != 0 {
    data = gin.H{
      "count": d.Count,
      "limit": d.Limit,
      "skip":  d.Skip,
    }
  }

  var list []gin.H
  for _, item := range d.List {
    list = append(list, fn(item))
  }

  data["list"] = list

  return data
}

func (d *DailyList) Ids() []bson.ObjectId {
  var list []bson.ObjectId
  for _, item := range d.List {
    list = append(list, item.Id)
  }
  return list
}
