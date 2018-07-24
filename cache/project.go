package cache

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/model"
)

// redis清项目信息
func ProjectDel(r redis.Conn, id string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  _, err := r.Do("DEL", n)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] ∆∆ Del |", n)
  }

  return err
}

// redis存项目信息
func ProjectSet(r redis.Conn, project *model.Project) {
  if project == nil {
    return
  }

  // 只存基本信息，不存关联表的信息
  m := project.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, project.Id.Hex())
  r.Do("SET", n, bytes)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Set |", n)
  }
}

// redis获取项目信息
func ProjectGet(r redis.Conn, id string, project *model.Project) bool {
  return false

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  data, _ := r.Do("GET", n)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    if gin.Mode() == gin.DebugMode {
      fmt.Println("[RDS] ∆∆ Get |", n)
    }
    return false
  }

  pid := res.Get("id").String()

  if !bson.IsObjectIdHex(pid) {
    if gin.Mode() == gin.DebugMode {
      fmt.Println("[RDS] ∆∆ Get |", n)
    }
    return false
  }

  project.Id = bson.ObjectIdHex(pid)
  project.Name = res.Get("name").String()
  project.Deadline = res.Get("deadline").Time()
  project.Description = res.Get("description").String()
  project.Weight = int(res.Get("weight").Int())
  project.Status = int(res.Get("status").Int())
  project.CreateTime = res.Get("createTime").Time()
  project.Exist = res.Get("exist").Bool()

  departments := res.Get("departments")
  if departments.IsArray() {
    for _, item := range departments.Array() {
      if item.Exists() {
        id := item.Get("id")
        if bson.IsObjectIdHex(id.String()) {
          project.Departments = append(project.Departments, mgo.DBRef{
            Collection: model.DepartmentCollection,
            Database:   conf.MgoDBName,
            Id:         bson.ObjectIdHex(id.String()),
          })
        }
      }
    }
  }

  missions := res.Get("missions")
  if missions.IsArray() {
    for _, item := range missions.Array() {
      if item.Exists() {
        id := item.Get("id")
        if bson.IsObjectIdHex(id.String()) {
          project.Missions = append(project.Missions, mgo.DBRef{
            Collection: model.MissionCollection,
            Database:   conf.MgoDBName,
            Id:         bson.ObjectIdHex(id.String()),
          })
        }
      }
    }
  }

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Get |", n)
  }

  return true
}
