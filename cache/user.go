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
func UserDel(id string) error {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, id)
  _, err := rd.Do("DEL", n)

  return err
}

// redis存用户信息
func UserSet(id string, user *model.User) {
  rd := db.RedisPool.Get()
  defer rd.Close()

  // 只存基本信息，不存关联表的信息
  m := user.GetMap(nil)
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, id)
  rd.Do("SET", n, bytes)
}

// redis获取用户信息
func UserGet(id string, user *model.User) bool {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.UserCollection, id)
  data, _ := rd.Do("GET", n)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    return false
  }

  uid := res.Get("id").String()

  if !bson.IsObjectIdHex(uid) {
    return false
  }

  user.Id = bson.ObjectIdHex(uid)
  user.CreateTime = res.Get("createTime").Time()
  user.Email = res.Get("email").String()
  user.NickName = res.Get("nickname").String()
  user.Role = int(res.Get("role").Int())
  user.Status = int(res.Get("status").Int())
  user.Title = res.Get("title").String()
  user.NickName = res.Get("username").String()
  user.Exist = true

  departmentId := res.Get("department.id").String()
  if bson.IsObjectIdHex(departmentId) {
    user.Department = mgo.DBRef{
      Collection: model.DepartmentCollection,
      Database:   conf.MgoDBName,
      Id:         bson.ObjectIdHex(departmentId),
    }
  }

  return true
}
