package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func GetProjectOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // check
  ctx.ErrorIfStringNotObjectId(id, errno.ErrProjectIdError)

  if ctx.HandleErrorIf() {
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
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrSkipRange)
  ctx.ErrorIfIntLessThen(limit, 1, errno.ErrLimitRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrLimitRange)

  // query
  data, err := service.GetProjectsList(skip, limit, nil)

  // check
  if err != nil {
    ctx.Error(errno.ErrProjectNotFound)
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
  ctx.ErrorIfStringIsEmpty(name, errno.ErrProjectNameEmpty)
  ctx.ErrorIfLenLessThen(name, 4, errno.ErrProjectNameTooShort)
  ctx.ErrorIfLenMoreThen(name, 15, errno.ErrProjectNameTooLong)
  ctx.ErrorIfTimeEarlierThen(deadline, time.Now(), errno.ErrProjectDeadlineTooSoon)
  ctx.ErrorIfIntIsZero(len(departments), errno.ErrProjectDepartmentsEmpty)
  if weight != 1 && weight != 2 && weight != 3 {
    ctx.Error(errno.ErrProjectWeightError)
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
      ctx.ErrorIfStringNotObjectId(department.Str, errno.ErrProjectDepartmentNotFound)
    }
  }

  if ctx.HandleErrorIf() {
    return
  }

  // create
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Departments: departmentsRef,
    Description: description,
    Weight:      weight,
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
  ctx.ErrorIfStringNotObjectId(id, errno.ErrProjectIdError)
  ctx.ErrorIfStringIsEmpty(name, errno.ErrProjectNameEmpty)
  ctx.ErrorIfLenLessThen(name, 4, errno.ErrProjectNameTooShort)
  ctx.ErrorIfLenMoreThen(name, 15, errno.ErrProjectNameTooLong)
  ctx.ErrorIfTimeEarlierThen(deadline, time.Now(), errno.ErrProjectDeadlineTooSoon)
  ctx.ErrorIfIntIsZero(len(departments), errno.ErrProjectDepartmentsEmpty)
  if weight != 1 && weight != 2 && weight != 3 {
    ctx.Error(errno.ErrProjectWeightError)
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
      ctx.ErrorIfStringNotObjectId(department.Str, errno.ErrProjectDepartmentNotFound)
    }
  }

  if ctx.HandleErrorIf() {
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
