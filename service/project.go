package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/errgo"
  "workerbook/model"
)

// 创建项目
func CreateProject(data model.Project, departments []gjson.Result) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrProjectNameEmpty)
  errgo.ErrorIfLenLessThen(data.Name, 4, errgo.ErrProjectNameTooShort)
  errgo.ErrorIfLenMoreThen(data.Name, 15, errgo.ErrProjectNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  errgo.ErrorIfIntIsZero(len(departments), errgo.ErrProjectDepartmentsEmpty)
  if data.Weight != 1 && data.Weight != 2 && data.Weight != 3 {
    return errors.New(errgo.ErrProjectWeightError)
  }

  // handle departments and check is really an objectId
  var departmentsRef []mgo.DBRef
  for _, department := range departments {
    if bson.IsObjectIdHex(department.Str) {
      departmentsRef = append(departmentsRef, mgo.DBRef{
        Collection: model.DepartmentCollection,
        Database:   conf.MgoDBName,
        Id:         bson.ObjectIdHex(department.Str),
      })
    } else {
      errgo.ErrorIfStringNotObjectId(department.Str, errgo.ErrProjectDepartmentNotFound)
    }
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  data.Departments = departmentsRef

  // insert it.
  err = mg.C(model.ProjectCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateProjectFailed)
  }

  return nil
}

// 更新项目
func UpdateProject(id string, data bson.M) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if name, ok := data["name"]; ok {
    errgo.ErrorIfLenLessThen(name.(string), 4, errgo.ErrProjectNameTooShort)
    errgo.ErrorIfLenMoreThen(name.(string), 15, errgo.ErrProjectNameTooLong)
  }

  if deadline, ok := data["deadline"]; ok {
    errgo.ErrorIfTimeEarlierThen(deadline.(time.Time), time.Now(), errgo.ErrProjectDeadlineTooSoon)
  }

  if weight, ok := data["weight"]; ok {
    weight := weight.(int)
    if weight != 1 && weight != 2 && weight != 3 {
      return errors.New(errgo.ErrProjectWeightError)
    }
  }

  // handle departments and check is really an objectId
  if departments, ok := data["departments"]; ok {
    departments := departments.([]gjson.Result)
    var departmentsRef []mgo.DBRef
    for _, department := range departments {
      if bson.IsObjectIdHex(department.Str) {
        departmentsRef = append(departmentsRef, mgo.DBRef{
          Collection: model.DepartmentCollection,
          Database:   conf.MgoDBName,
          Id:         bson.ObjectIdHex(department.Str),
        })
      } else {
        errgo.ErrorIfStringNotObjectId(department.Str, errgo.ErrProjectDepartmentNotFound)
      }
    }
    data["departments"] = departmentsRef
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 先要清缓存，清成功后才可以更新数据
  err = cache.ProjectDel(id)

  if err != nil {
    return errors.New(errgo.ErrUpdateProjectFailed)
  }

  // update
  data["exist"] = true
  err = mg.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateProjectFailed)
  }

  return nil
}

// 根据id查找项目
func GetProjectInfoById(id string, refs ... string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Project)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.ProjectGet(id, data)
  if !rok {
    err = mg.C(model.ProjectCollection).Find(bson.M{
      "_id":   bson.ObjectIdHex(id),
      "exist": true,
    }).One(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrProjectNotFound)
    }
    return nil, err
  }

  // 存到缓存
  if !rok {
    cache.ProjectSet(id, data)
  }

  return data.GetMap(mg, refs...), nil
}

// 根据id删除项目
func DelProjectById(id string) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 清除缓存，缓存清成功才可以清数据，不然会有脏数据
  err = cache.ProjectDel(id)

  if err != nil {
    return errors.New(errgo.ErrDeleteProjectFailed)
  }

  // 删除数据
  err = mg.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": bson.M{
      "exist": false,
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrProjectNotFound)
    }
    return err
  }

  if err != nil {
    return errors.New(errgo.ErrDeleteProjectFailed)
  }

  return nil
}

// 查找项目列表(当skip和limit都为0时，查找全部)
func GetProjectsList(skip int, limit int, query bson.M, refs ... string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  if skip != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.Project)
  query["exist"] = true

  // find it
  if skip == 0 && limit == 0 {
    err = mg.C(model.ProjectCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = mg.C(model.ProjectCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrProjectNotFound)
    }
    return nil, err
  }

  // result
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(mg, refs...))
  }

  if skip == 0 && limit == 0 {
    return gin.H{
      "list": list,
    }, nil
  }

  // get count
  count, err := mg.C(model.ProjectCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrProjectNotFound)
  }

  return gin.H{
    "list":  list,
    "count": count,
    "skip":  skip,
    "limit": limit,
  }, nil
}
