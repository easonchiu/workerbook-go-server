package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  `gopkg.in/mgo.v2/bson`
  `time`
  "workerbook/errno"
  `workerbook/model`
  "workerbook/mongo"
)

// 创建项目
func CreateProject(data model.Project) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // supplement other data.
  data.CreateTime = time.Now()
  data.Status = 1

  // insert it.
  err = db.C(model.ProjectCollection).Insert(data)

  if err != nil {
    return errors.New(errno.ErrCreateProjectFailed)
  }

  return nil
}

// 更新项目
func UpdateProject(id bson.ObjectId, data model.Project) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // update
  err = db.C(model.ProjectCollection).UpdateId(id, bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errno.ErrUpdateProjectFailed)
  }

  return nil
}

// 根据id查找项目
func GetProjectInfoById(id bson.ObjectId) (gin.H, error) {
  db, close, err := mongo.CloneDB()

  data := new(model.Project)

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  err = db.C(model.ProjectCollection).FindId(id).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errno.ErrProjectNotFound)
    }
    return nil, err
  }

  return data.GetMap(db), nil
}

// 查找项目列表
// 当skip和limit都为0时，查找全部
func GetProjectsList(skip int, limit int, query bson.M) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  data := new([]model.Project)

  if limit < 0 {
    limit = 0
  } else if limit > 100 {
    limit = 100
  }

  // find it
  if skip == 0 && limit == 0 {
    err = db.C(model.ProjectCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = db.C(model.ProjectCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errno.ErrProjectNotFound)
    }
    return nil, err
  }

  // get count
  count, err := db.C(model.ProjectCollection).Count()

  if err != nil {
    return nil, errors.New(errno.ErrProjectNotFound)
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

  return gin.H{
    "list":  list,
    "count": count,
    "skip":  skip,
    "limit": limit,
  }, nil
}
