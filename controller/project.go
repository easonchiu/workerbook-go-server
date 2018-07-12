package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func GetProjectOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // query
  projectInfo, err := service.GetProjectInfoById(bson.ObjectIdHex(id))

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": projectInfo,
  })
}

// 获取项目列表
func GetProjectsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
  errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
  errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)

  // query
  data, err := service.GetProjectsList(skip, limit, nil)

  // check
  if err != nil {
    ctx.Error(errgo.ErrProjectNotFound)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data,
  })
}

// 创建项目
func CreateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  name := ctx.getRaw("name")
  deadline := ctx.getRawTime("deadline")
  departments := ctx.getRawArray("departments")
  description := ctx.getRaw("description")
  weight := ctx.getRawInt("weight")

  // check
  errgo.ErrorIfStringIsEmpty(name, errgo.ErrProjectNameEmpty)
  errgo.ErrorIfLenLessThen(name, 4, errgo.ErrProjectNameTooShort)
  errgo.ErrorIfLenMoreThen(name, 15, errgo.ErrProjectNameTooLong)
  errgo.ErrorIfTimeEarlierThen(deadline, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  errgo.ErrorIfIntIsZero(len(departments), errgo.ErrProjectDepartmentsEmpty)
  if weight != 1 && weight != 2 && weight != 3 {
    ctx.Error(errgo.ErrProjectWeightError)
    return
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

  if errgo.HandleError(ctx.Error) {
    return
  }

  // create
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Departments: departmentsRef,
    Description: description,
    Weight:      weight,
    CreateTime:  time.Now(),
    Status:      1,
  }

  // insert
  err := service.CreateProject(data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改项目
func UpdateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  name := ctx.getRaw("name")
  deadline := ctx.getRawTime("deadline")
  departments := ctx.getRawArray("departments")
  description := ctx.getRaw("description")
  weight := ctx.getRawInt("weight")

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)
  errgo.ErrorIfStringIsEmpty(name, errgo.ErrProjectNameEmpty)
  errgo.ErrorIfLenLessThen(name, 4, errgo.ErrProjectNameTooShort)
  errgo.ErrorIfLenMoreThen(name, 15, errgo.ErrProjectNameTooLong)
  errgo.ErrorIfTimeEarlierThen(deadline, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  errgo.ErrorIfIntIsZero(len(departments), errgo.ErrProjectDepartmentsEmpty)
  if weight != 1 && weight != 2 && weight != 3 {
    ctx.Error(errgo.ErrProjectWeightError)
    return
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

  if errgo.HandleError(ctx.Error) {
    return
  }

  // update
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Departments: departmentsRef,
    Description: description,
    Weight:      weight,
  }

  err := service.UpdateProject(bson.ObjectIdHex(id), data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
