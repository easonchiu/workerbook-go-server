package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/mongo"
)

// 创建部门
func CreateDepartment(data model.Department) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // 设置缺省值
  data.CreateTime = time.Now()
  data.UserCount = 0

  // 名称不能重复
  count, err := db.C(model.DepartmentCollection).Find(bson.M{"name": data.Name}).Count()

  if err != nil {
    return err
  }

  if count > 0 {
    return errors.New(errno.ErrSameDepartmentName)
  }

  // insert it.
  err = db.C(model.DepartmentCollection).Insert(data)

  if err != nil {
    return errors.New(errno.ErrCreateDepartmentFailed)
  }

  return nil
}

// Query group info by id.
func GetDepartmentInfoById(id bson.ObjectId) (*model.Department, error) {
  db, close, err := mongo.CloneDB()

  data := new(model.Department)

  if err != nil {
    return data, err
  } else {
    defer close()
  }

  err = db.C(model.DepartmentCollection).FindId(id).One(data)

  if err != nil {
    return data, err
  }

  return data, nil
}

// Query groups list with skip and limit.
func GetDepartmentsList(skip int, limit int) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  data := new([]model.Department)

  if limit < 0 {
    limit = 0
  } else if limit > 100 {
    limit = 100
  }

  // find it
  err = db.C(model.DepartmentCollection).Find(nil).Skip(skip).Limit(limit).Sort("-createTime").All(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errno.ErrDepartmentNotFound)
    }
    return nil, err
  }

  // get count
  count, err := db.C(model.DepartmentCollection).Count()

  if err != nil {
    return nil, errors.New(errno.ErrDepartmentNotFound)
  }

  // result
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(db))
  }

  return gin.H{
    "list":  list,
    "count": count,
    "skip":  skip,
    "limit": limit,
  }, nil
}

// 全量更新所有部门的人数
func UpdateDepartmentsUserCount() error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  departments := new([]model.Department)
  err = db.C(model.DepartmentCollection).Find(nil).All(departments)

  if err != nil {
    return err
  }

  for _, department := range *departments {
    count, err := db.C(model.UserCollection).Find(bson.M{"department.$id": department.Id}).Count()
    if err != nil {
      return err
    }
    db.C(model.DepartmentCollection).UpdateId(department.Id, bson.M{
      "$set": bson.M{
        "userCount": count,
      },
    })
  }

  return nil
}

func UpdateDepartment(id bson.ObjectId, m bson.M) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 名称唯一
  count, err := db.C(model.DepartmentCollection).Find(bson.M{
    "name": m["name"],
    "_id": bson.M{
      "$ne": id,
    },
  }).Count()

  if err != nil {
    return errors.New(errno.ErrUpdateDepartmentFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameDepartmentName)
  }

  // 更新数据
  err = db.C(model.DepartmentCollection).UpdateId(id, bson.M{
    "$set": m,
  })

  if err != nil {
    return errors.New(errno.ErrUpdateDepartmentFailed)
  }

  return nil
}
