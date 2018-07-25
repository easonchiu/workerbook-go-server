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

// redis清任务信息
func MissionDel(r redis.Conn, id string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, id)
  _, err := r.Do("DEL", n)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] ∆∆ Del |", n)
  }

  return err
}

// redis存任务信息
func MissionSet(r redis.Conn, mission *model.Mission) {
  if mission == nil {
    return
  }

  m := mission.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, mission.Id.Hex())
  r.Do("SET", n, bytes)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Set |", n)
  }
}

// redis获取任务信息
func MissionGet(r redis.Conn, id string, mission *model.Mission) bool {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, id)
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

  mid := res.Get("id").String()

  if !bson.IsObjectIdHex(mid) {
    if gin.Mode() == gin.DebugMode {
      fmt.Println("[RDS] ∆∆ Get |", n)
    }
    return false
  }

  mission.Id = bson.ObjectIdHex(mid)
  mission.Name = res.Get("name").String()

  mission.Progress = int(res.Get("progress").Int())
  mission.Deadline = res.Get("deadline").Time()
  mission.Status = int(res.Get("status").Int())

  projectId := res.Get("projectId").String()
  if bson.IsObjectIdHex(projectId) {
    mission.ProjectId = bson.ObjectIdHex(projectId)
  }

  userId := res.Get("user.id").String()
  if bson.IsObjectIdHex(userId) {
    mission.User = mgo.DBRef{
      Id:         bson.ObjectIdHex(userId),
      Collection: model.UserCollection,
      Database:   conf.MgoDBName,
    }
  }

  mission.CreateTime = res.Get("createTime").Time()
  mission.Exist = res.Get("exist").Bool()

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Get |", n)
  }

  return true
}
