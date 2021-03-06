package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// collection name
const MissionCollection = "missions"

// collection schema
type Mission struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 任务名
  Name string `bson:"name"`

  // 之前的进度
  PreProgress int `bson:"preProgress"`

  // 统计时间(2006-01-02格式)
  ChartTime string `bson:"chartTime"`

  // 进度
  Progress int `bson:"progress"`

  // 执行人
  User mgo.DBRef `bson:"user"`

  // 截至时间
  Deadline time.Time `bson:"deadline"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`

  // 所属任务
  Project mgo.DBRef `bson:"project"`

  // 是否存在
  Exist bool `bson:"exist"`

  // 操作人
  Editor mgo.DBRef `bson:"editor,omitempty"`

  // 操作时间
  EditTime time.Time `bson:"editTime,omitempty"`
}

func (m Mission) GetMap(forgets ... string) gin.H {
  data := gin.H{
    "id":         m.Id,
    "name":       m.Name,
    "createTime": m.CreateTime,
    "chartTime":  m.ChartTime,
    "deadline":   m.Deadline,
    "isTimeout":  m.Deadline.Before(time.Now()),
    "project": gin.H{
      "id": m.Project.Id,
    },
    "preProgress": m.PreProgress,
    "progress":    m.Progress,
    "user": gin.H{
      "id": m.User.Id,
    },
    "exist": m.Exist,
    "editor": bson.M{
      "id": m.Editor.Id,
    },
    "editTime": m.EditTime,
  }

  if forgets != nil {
    if forgets[0] == REMEMBER {
      remember(data, forgets[1:]...)
    } else {
      forget(data, forgets...)
    }
  }

  return data
}

// 任务列表结构
type MissionList struct {
  List  []*Mission
  Count int
  Limit int
  Skip  int
}

// 列表的迭代器
func (d *MissionList) Each(fn func(*Mission) gin.H) gin.H {
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

func (d *MissionList) Find(id bson.ObjectId) *Mission {
  if d.List == nil {
    return nil
  }
  for _, item := range d.List {
    if item.Id == id {
      return item
    }
  }
  return nil
}

func (d *MissionList) Ids() []bson.ObjectId {
  var list []bson.ObjectId
  for _, item := range d.List {
    list = append(list, item.Id)
  }
  return list
}
