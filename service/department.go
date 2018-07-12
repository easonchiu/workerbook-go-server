package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/errgo"
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
    return errors.New(errgo.ErrSameDepartmentName)
  }

  // insert it.
  err = db.C(model.DepartmentCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateDepartmentFailed)
  }

  return nil
}

// 根据id查找部门信息
func GetDepartmentInfoById(id bson.ObjectId) (gin.H, error) {
  db, close, err := mongo.CloneDB()

  data := new(model.Department)

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  err = db.C(model.DepartmentCollection).FindId(id).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrDepartmentNotFound)
    }
    return nil, err
  }

  return data.GetMap(db), nil
}

// 查找部门列表
// 当skip和limit都为0时，查找全部
func GetDepartmentsList(skip int, limit int, query bson.M) (gin.H, error) {
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
  if skip == 0 && limit == 0 {
    err = db.C(model.DepartmentCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = db.C(model.DepartmentCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrDepartmentNotFound)
    }
    return nil, err
  }

  // result
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(db))
  }

  if skip == 0 && limit == 0 {
    return gin.H{
      "list":  list,
    }, nil
  }

  // get count
  count, err := db.C(model.DepartmentCollection).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrDepartmentNotFound)
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

func UpdateDepartment(id bson.ObjectId, data model.Department) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 名称唯一
  count, err := db.C(model.DepartmentCollection).Find(bson.M{
    "name": data.Name,
    "_id": bson.M{
      "$ne": id,
    },
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameDepartmentName)
  }

  // 更新数据
  err = db.C(model.DepartmentCollection).UpdateId(id, bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  return nil
}
