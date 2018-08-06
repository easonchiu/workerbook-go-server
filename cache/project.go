package cache

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "reflect"
  "strconv"
  "workerbook/conf"
  "workerbook/model"
)

// redis清项目信息
func ProjectDel(r redis.Conn, id string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  _, err := r.Do("DEL", n)

  if gin.IsDebugging() {
    fmt.Println("[RDS] 🗑 Del |", n)
  }

  return err
}

// redis存项目信息
func ProjectSet(r redis.Conn, project *model.Project) {
  if project == nil {
    return
  }

  m := project.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, project.Id.Hex())
  if conf.RedisExpireTime != 0 {
    r.Do("SETEX", n, conf.RedisExpireTime, bytes)
  } else {
    r.Do("SET", n, bytes)
  }

  if gin.IsDebugging() {
    fmt.Println("[RDS] ✨ Set |", n)
  }
}

// redis获取项目信息
func ProjectGet(r redis.Conn, id string, project *model.Project) bool {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  data, _ := r.Do("GET", n)

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

  pid := res.Get("id").String()

  if !bson.IsObjectIdHex(pid) {
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️ Get |", n)
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
            Database:   conf.MgoDBName,
            Collection: model.DepartmentCollection,
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
            Database:   conf.MgoDBName,
            Collection: model.MissionCollection,
            Id:         bson.ObjectIdHex(id.String()),
          })
        }
      }
    }
  }

  editorId := res.Get("editor.id").String()
  if bson.IsObjectIdHex(editorId) {
    project.Editor = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         bson.ObjectIdHex(editorId),
    }
  }
  project.EditTime = res.Get("editTime").Time()

  if gin.IsDebugging() {
    fmt.Println("[RDS] ⚡️ Get |", n)
  }

  return true
}

// redis存项目进度
func ProjectProgressSet(r redis.Conn, id string, progress int) {

  n := fmt.Sprintf("%v:%v:%v:%v", conf.MgoDBName, model.ProjectCollection, "progress", id)

  if conf.RedisExpireTime != 0 {
    r.Do("SETEX", n, conf.RedisExpireTime, progress)
  } else {
    r.Do("SET", n, progress)
  }

  if gin.IsDebugging() {
    fmt.Println("[RDS] ✨ Set |", n)
  }
}

// redis获取项目进度
func ProjectProgressGet(r redis.Conn, id string) (int, bool) {

  n := fmt.Sprintf("%v:%v:%v:%v", conf.MgoDBName, model.ProjectCollection, "progress", id)

  data, _ := r.Do("GET", n)

  t := reflect.TypeOf(data)

  if t == nil {
    return 0, false
  }

  // 如果是整数类型，返回该值
  if t.String() == "[]uint8" {
    str := string(data.([]uint8))
    i, err := strconv.Atoi(str)
    return i, err == nil
  }

  return 0, false
}

// redis删除项目进度
func ProjectProgressDel(r redis.Conn, id string) error {

  n := fmt.Sprintf("%v:%v:%v:%v", conf.MgoDBName, model.ProjectCollection, "progress", id)

  _, err := r.Do("DEL", n)

  return err
}