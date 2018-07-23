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

// redis清用户信息
func ProjectDel(id string) error {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  _, err := rd.Do("DEL", n)

  return err
}

// redis存用户信息
func ProjectSet(id string, project *model.Project) {
  rd := db.RedisPool.Get()
  defer rd.Close()

  // 只存基本信息，不存关联表的信息
  m := project.GetMap(nil)
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  rd.Do("SET", n, bytes)
}

// redis获取用户信息
func ProjectGet(id string, project *model.Project) bool {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.ProjectCollection, id)
  data, _ := rd.Do("GET", n)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    return false
  }

  pid := res.Get("id").String()

  if !bson.IsObjectIdHex(pid) {
    return false
  }

  project.Id = bson.ObjectIdHex(pid)
  project.Name = res.Get("name").String()
  project.Deadline = res.Get("deadline").Time()
  project.Description = res.Get("description").String()
  project.Weight = int(res.Get("weight").Int())
  project.Status = int(res.Get("status").Int())
  project.CreateTime = res.Get("createTime").Time()
  project.Exist = true

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

  return true
}
