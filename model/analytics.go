package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "sort"
  "strconv"
  "time"
)

// collection name
const AnalyticsCollection = "analytics"

// collection schema
type Analytics struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 用户名
  UserId bson.ObjectId `bson:"userId"`

  // department
  DepartmentId bson.ObjectId `bson:"departmentId"`

  // 任务id
  MissionId bson.ObjectId `bson:"missionId"`

  // 进度
  Progress int `bson:"progress"`

  // 项目id
  ProjectId bson.ObjectId `bson:"projectId"`

  // 日期(2018-02-01这样的格式)
  Day string `bson:"day"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`
}

// 单部门数据
type DepartmentAnalytics struct {
  Department *Department
  Missions   []*Mission
}

// 部门列表数据
type DepartmentAnalyticsList struct {
  List  []*DepartmentAnalytics
  Skip  int
  Limit int
  Count int
}

// 部门列表中找到指定部门
func (d *DepartmentAnalyticsList) Find(id bson.ObjectId) (department *DepartmentAnalytics) {
  for _, item := range d.List {
    if item.Department.Id == id {
      return item
    }
  }
  return nil
}

// 单项目数据
type ProjectAnalytics struct {
  Project  *Project
  Missions []*Mission
}

// 项目列表数据
type ProjectListAnalytics struct {
  List  []*ProjectAnalytics
  Skip  int
  Limit int
  Count int
}

// 单用户数据
type UserAnalytics struct {
  User     *User
  Missions []*Mission
}

// 部门数据
type DepartmentUsersAnalytics struct {
  Department *Department
  Users      []*UserAnalytics
}

func (d *DepartmentUsersAnalytics) Find(id bson.ObjectId) (user *UserAnalytics) {
  for _, item := range d.Users {
    if item.User.Id == id {
      return item
    }
  }
  return nil
}

// 任务单天的数据
type MissionChartData struct {
  Progress  int
  Day       string
}

// 任务数据
type MissionChartAnalytics struct {
  Id        bson.ObjectId
  Name      string
  IsTimeout bool
  Deadline  time.Time
  Data      []*MissionChartData
}

// 往任务中添加单天的数据（天去重）
func (m *MissionChartAnalytics) Append(data *MissionChartData) {
  exist := false
  for _, item := range m.Data {
    if item.Day == data.Day {
      exist = true
      break
    }
  }
  if !exist {
    m.Data = append(m.Data, data)
  }
}

// 获取任务数据的map
func (p *MissionChartAnalytics) GetMap() gin.H {
  data := gin.H{
    "name":      p.Name,
    "id":        p.Id,
    "isTimeout": p.IsTimeout,
    "deadline":  p.Deadline,
  }

  chartData := make([]gin.H, 0, 120)

  for _, item := range p.Data {
    chartData = append(chartData, gin.H{
      "day":       item.Day,
      "progress":  item.Progress,
    })
  }

  sort.Slice(chartData, func(i, j int) bool {
    d1, ok := chartData[i]["day"]
    if !ok {
      return true
    }

    d2, ok := chartData[j]["day"]
    if !ok {
      return true
    }

    sd1, err := strconv.Atoi(d1.(string))
    if err != nil {
      return true
    }

    sd2, err := strconv.Atoi(d2.(string))
    if err != nil {
      return true
    }

    return sd1 < sd2
  })

  data["data"] = chartData

  return data
}

// 项目的数据
type ProjectChartAnalytics struct {
  Project  *Project
  Missions []*MissionChartAnalytics
}

// 项目数据中找到某个任务
func (p *ProjectChartAnalytics) Find(id bson.ObjectId) (res *MissionChartAnalytics) {
  for _, item := range p.Missions {
    if item.Id == id {
      return item
    }
  }
  return nil
}

// 获取项目数据map
func (p *ProjectChartAnalytics) GetMap() gin.H {
  data := p.Project.GetMap(REMEMBER, "id", "name", "deadline", "createTime", "isTimeout")

  missions := make([]gin.H, 0, len(p.Missions))
  for _, item := range p.Missions {
    missions = append(missions, item.GetMap())
  }
  data["missions"] = missions

  return data
}
