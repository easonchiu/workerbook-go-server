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

  if gin.IsDebugging() {
    fmt.Println("[RDS] 🗑 Del |", n)
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
  if conf.RedisExpireTime != 0 {
    r.Do("SETEX", n, conf.RedisExpireTime, bytes)
  } else {
    r.Do("SET", n, bytes)
  }

  if gin.IsDebugging() {
    fmt.Println("[RDS] ✨ Set |", n)
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
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️️ Get |", n)
    }
    return false
  }

  mid := res.Get("id").String()

  if !bson.IsObjectIdHex(mid) {
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️ Get |", n)
    }
    return false
  }

  mission.Id = bson.ObjectIdHex(mid)
  mission.Name = res.Get("name").String()
  mission.PreProgress = int(res.Get("preProgress").Int())
  mission.Progress = int(res.Get("progress").Int())
  mission.Deadline = res.Get("deadline").Time()
  mission.ChartTime = res.Get("chartTime").String()

  projectId := res.Get("project.id").String()
  if bson.IsObjectIdHex(projectId) {
    mission.Project = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.ProjectCollection,
      Id:         bson.ObjectIdHex(projectId),
    }
  }

  userId := res.Get("user.id").String()
  if bson.IsObjectIdHex(userId) {
    mission.User = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         bson.ObjectIdHex(userId),
    }
  }

  mission.CreateTime = res.Get("createTime").Time()
  mission.Exist = res.Get("exist").Bool()

  editorId := res.Get("editor.id").String()
  if bson.IsObjectIdHex(editorId) {
    mission.Editor = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         bson.ObjectIdHex(editorId),
    }
  }
  mission.EditTime = res.Get("editTime").Time()

  if gin.IsDebugging() {
    fmt.Println("[RDS] ⚡️ Get |", n)
  }

  return true
}
