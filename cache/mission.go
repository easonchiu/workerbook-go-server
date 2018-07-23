package cache

import (
  "encoding/json"
  "fmt"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/model"
)

// redis清任务信息
func MissionDel(id string) error {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, id)
  _, err := rd.Do("DEL", n)

  return err
}

// redis存任务信息
func MissionSet(id string, mission *model.Mission) {
  rd := db.RedisPool.Get()
  defer rd.Close()

  // 只存基本信息，不存关联表的信息
  m := mission.GetMap(nil)
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, id)
  rd.Do("SET", n, bytes)
}

// redis获取任务信息
func MissionGet(id string, mission *model.Mission) bool {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.MissionCollection, id)
  data, _ := rd.Do("GET", n)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    return false
  }

  mid := res.Get("id").String()

  if !bson.IsObjectIdHex(mid) {
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
      Collection: model.MissionCollection,
      Database:   conf.MgoDBName,
    }
  }

  mission.CreateTime = res.Get("createTime").Time()
  mission.Exist = true

  return true
}
