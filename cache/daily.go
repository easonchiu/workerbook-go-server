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

// redis清日报信息
func DailyDel(r redis.Conn, uid string, day string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DailyCollection, day)
  _, err := r.Do("HDEL", n, uid)

  if gin.IsDebugging() {
    fmt.Println("[RDS] 🗑 Del |", n)
  }

  return err
}

// redis存日报信息
func DailySet(r redis.Conn, uid string, day string, daily *model.Daily) {
  if daily == nil {
    return
  }

  m := daily.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DailyCollection, day)
  r.Do("HSET", n, uid, bytes)

  if gin.IsDebugging() {
    fmt.Println("[RDS] ✨ Set |", n)
  }
}

// redis获取日报信息
func DailyGet(r redis.Conn, uid string, day string, daily *model.Daily) bool {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DailyCollection, day)
  data, _ := r.Do("HGET", n, uid)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️ Get |", n)
    }
    return false
  }

  did := res.Get("id").String()

  if !bson.IsObjectIdHex(did) {
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️ Get |", n)
    }
    return false
  }

  daily.Id = bson.ObjectIdHex(did)
  daily.Day = res.Get("day").String()
  daily.CreateTime = res.Get("createTime").Time()
  daily.UpdateTime = res.Get("updateTime").Time()
  daily.DepartmentName = res.Get("departmentName").String()

  var dailies []model.DailyItem
  if res.Get("dailies").IsArray() {
    arr := res.Get("dailies").Array()
    for _, item := range arr {
      var di = new(model.DailyItem)
      id := item.Get("id").String()
      content := item.Get("content").String()
      progress := int(item.Get("progress").Int())
      missionName := item.Get("mission.name").String()
      missionId := item.Get("mission.id").String()
      projectName := item.Get("project.name").String()
      projectId := item.Get("project.id").String()

      if bson.IsObjectIdHex(id) && bson.IsObjectIdHex(missionId) && bson.IsObjectIdHex(projectId) {
        di.Id = bson.ObjectIdHex(id)
        di.Content = content
        di.Progress = progress
        di.MissionName = missionName
        di.MissionId = bson.ObjectIdHex(missionId)
        di.ProjectName = projectName
        di.ProjectId = bson.ObjectIdHex(projectId)

        dailies = append(dailies, *di)
      }
    }
  }

  daily.Dailies = dailies

  userId := res.Get("user.id").String()
  if bson.IsObjectIdHex(userId) {
    daily.User = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         bson.ObjectIdHex(userId),
    }
  }

  fmt.Println("redis daily", daily)

  if gin.IsDebugging() {
    fmt.Println("[RDS] ⚡️ Get |", n)
  }

  return true
}
