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

// redis清用户信息
func UserDel(r redis.Conn, id string) error {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, id)
  _, err := r.Do("DEL", n)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] ∆∆ Del |", n)
  }

  return err
}

// redis存用户信息
func UserSet(r redis.Conn, user *model.User) {
  if user == nil {
    return
  }

  m := user.GetMap()
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, user.Id.Hex())
  r.Do("SET", n, bytes)

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Set |", n)
  }
}

// redis获取用户信息
func UserGet(r redis.Conn, id string, user *model.User) bool {
  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, id)
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

  uid := res.Get("id").String()

  if !bson.IsObjectIdHex(uid) {
    if gin.Mode() == gin.DebugMode {
      fmt.Println("[RDS] ∆∆ Get |", n)
    }
    return false
  }

  user.Id = bson.ObjectIdHex(uid)
  user.CreateTime = res.Get("createTime").Time()
  user.Email = res.Get("email").String()
  user.NickName = res.Get("nickname").String()
  user.Role = int(res.Get("role").Int())
  user.Status = int(res.Get("status").Int())
  user.Title = res.Get("title").String()
  user.UserName = res.Get("username").String()
  user.Exist = res.Get("exist").Bool()

  departmentId := res.Get("department.id").String()
  if bson.IsObjectIdHex(departmentId) {
    user.Department = mgo.DBRef{
      Collection: model.DepartmentCollection,
      Database:   conf.MgoDBName,
      Id:         bson.ObjectIdHex(departmentId),
    }
  }

  if gin.Mode() == gin.DebugMode {
    fmt.Println("[RDS] √√ Get |", n)
  }

  return true
}
