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

// redis清部门信息
func DepartmentDel(r redis.Conn, id string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, id)
  _, err := r.Do("DEL", n)

  if gin.IsDebugging() {
    fmt.Println("[RDS] 🗑 Del |", n)
  }

  return err
}

// redis存部门信息
func DepartmentSet(r redis.Conn, department *model.Department) {
  if department == nil {
    return
  }

  m := department.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, department.Id.Hex())
  r.Do("SET", n, bytes)
  r.Do("EXPIRE", n, conf.RedisExpireTime)

  if gin.IsDebugging() {
    fmt.Println("[RDS] ✨ Set |", n)
  }
}

// redis获取部门信息
func DepartmentGet(r redis.Conn, id string, department *model.Department) bool {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, id)
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

  did := res.Get("id").String()

  if !bson.IsObjectIdHex(did) {
    if gin.IsDebugging() {
      fmt.Println("[RDS] ⚠️ Get |", n)
    }
    return false
  }

  department.Id = bson.ObjectIdHex(did)
  department.CreateTime = res.Get("createTime").Time()
  department.Name = res.Get("name").String()
  department.UserCount = int(res.Get("userCount").Int())
  department.Exist = res.Get("exist").Bool()

  editorId := res.Get("editor.id").String()
  if bson.IsObjectIdHex(editorId) {
    department.Editor = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         bson.ObjectIdHex(editorId),
    }
  }
  department.EditTime = res.Get("editTime").Time()

  if gin.IsDebugging() {
    fmt.Println("[RDS] ⚡️ Get |", n)
  }

  return true
}
