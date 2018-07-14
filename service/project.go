package service

import (
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/mongo"
)

// 创建项目
func CreateProject(data model.Project, departments []gjson.Result) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // 是否存在的标志
  data.Exist = true

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
        Database:   conf.DBName,
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
  err = db.C(model.ProjectCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateProjectFailed)
  }

  return nil
}

// 更新项目
func UpdateProject(id string, data model.Project, departments []gjson.Result) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrProjectNameEmpty)
  errgo.ErrorIfLenLessThen(data.Name, 4, errgo.ErrProjectNameTooShort)
  errgo.ErrorIfLenMoreThen(data.Name, 15, errgo.ErrProjectNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  if data.Weight != 1 && data.Weight != 2 && data.Weight != 3 {
    return errors.New(errgo.ErrProjectWeightError)
  }

  // handle departments and check is really an objectId
  if len(departments) > 0 {
    var departmentsRef []mgo.DBRef
    for _, department := range departments {
      if bson.IsObjectIdHex(department.Str) {
        departmentsRef = append(departmentsRef, mgo.DBRef{
          Collection: model.DepartmentCollection,
          Database:   conf.DBName,
          Id:         bson.ObjectIdHex(department.Str),
        })
      } else {
        errgo.ErrorIfStringNotObjectId(department.Str, errgo.ErrProjectDepartmentNotFound)
      }
    }
    data.Departments = departmentsRef
  } else {
    data.Departments = nil
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // update
  data.Exist = true
  err = db.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateProjectFailed)
  }

  return nil
}

// 根据id查找项目
func GetProjectInfoById(id string) (gin.H, error) {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Project)

  err = db.C(model.ProjectCollection).Find(model.Project{
    Id:    bson.ObjectIdHex(id),
    Exist: true,
  }).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrProjectNotFound)
    }
    return nil, err
  }

  return data.GetMap(db), nil
}

// 根据id删除项目
func DelProjectById(id string) error {
  db, closer, err := mongo.CloneDB()

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

  // 删除
  err = db.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": model.Project{
      Exist: false,
    },
  })

  if err != nil {
    fmt.Println(err)
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
func GetProjectsList(skip int, limit int, query model.Project) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
  errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
  errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.Project)
  query.Exist = true

  if limit < 0 {
    limit = 0
  } else if limit > 100 {
    limit = 100
  }

  fmt.Println(query)

  // find it
  if skip == 0 && limit == 0 {
    err = db.C(model.ProjectCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = db.C(model.ProjectCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
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
    list = append(list, r.GetMap(db))
  }

  if skip == 0 && limit == 0 {
    return gin.H{
      "list": list,
    }, nil
  }

  // get count
  count, err := db.C(model.ProjectCollection).Find(query).Count()

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
