package cache

import (
  "encoding/json"
  "fmt"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/model"
)

// redis清部门信息
func DepartmentDel(id string) error {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, id)
  _, err := rd.Do("DEL", n)

  return err
}

// redis存部门信息
func DepartmentSet(id string, department *model.Department) {
  rd := db.RedisPool.Get()
  defer rd.Close()

  // 只存基本信息，不存关联表的信息
  m := department.GetMap(nil)
  bytes, _ := json.Marshal(m)

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, id)
  rd.Do("SET", n, bytes)
}

// redis获取部门信息
func DepartmentGet(id string, department *model.Department) bool {
  rd := db.RedisPool.Get()
  defer rd.Close()

  n := fmt.Sprintf("%v:%v:%v", conf.MgoDBName, model.DepartmentCollection, id)
  data, _ := rd.Do("GET", n)

  if data == nil {
    return false
  }

  res := gjson.ParseBytes(data.([]byte))

  if !res.Exists() {
    return false
  }

  did := res.Get("id").String()

  if !bson.IsObjectIdHex(did) {
    return false
  }

  department.Id = bson.ObjectIdHex(did)
  department.CreateTime = res.Get("createTime").Time()
  department.Name = res.Get("name").String()
  department.UserCount = int(res.Get("userCount").Int())
  department.Exist = true

  return true
}
