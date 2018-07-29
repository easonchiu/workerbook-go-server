package service

import (
  "errors"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/util"
)

// 创建项目
func CreateProject(ctx *context.New, data model.Project, departments []gjson.Result) error {

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data.Editor = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.UserCollection,
    Id:         bson.ObjectIdHex(ownUserId),
  }
  data.EditTime = time.Now()

  // 修改项目结束时间到今天的最晚时间，即23:59:59
  {
    y1, m1, d1 := data.Deadline.Local().Date()
    data.Deadline = time.Date(y1, m1, d1, 23, 59, 59, 0, time.Local)
  }

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrProjectNameEmpty)
  errgo.ErrorIfLenLessThen(data.Name, 4, errgo.ErrProjectNameTooShort)
  errgo.ErrorIfLenMoreThen(data.Name, 15, errgo.ErrProjectNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  errgo.ErrorIfIntIsZero(len(departments), errgo.ErrProjectDepartmentsEmpty)
  if data.Weight != 1 && data.Weight != 2 && data.Weight != 3 {
    return errors.New(errgo.ErrProjectWeightError)
  }

  // 验证每个部门是否正常
  var departmentsRef []mgo.DBRef
  for _, item := range departments {
    itemId := item.String()
    if item.Exists() && bson.IsObjectIdHex(itemId) {
      _, err := GetDepartmentInfoById(ctx, itemId)
      if err == nil {
        departmentsRef = append(departmentsRef, mgo.DBRef{
          Database:   conf.MgoDBName,
          Collection: model.DepartmentCollection,
          Id:         bson.ObjectIdHex(itemId),
        })
      } else {
        return errors.New(errgo.ErrProjectDepartmentNotFound)
      }
    } else {
      return errors.New(errgo.ErrProjectDepartmentNotFound)
    }
  }

  data.Departments = departmentsRef

  // insert it.
  err := ctx.MgoDB.C(model.ProjectCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateProjectFailed)
  }

  return nil
}

// 更新项目
func UpdateProject(ctx *context.New, id string, data bson.M) error {

  if data == nil {
    return errors.New(errgo.ErrServerError)
  }

  // 限制更新字段
  util.Only(
    data,
    util.Keys{
      "name":        util.TypeString,
      "status":      util.TypeInt,
      "deadline":    util.TypeTime,
      "departments": util.TypeAny,
      "description": util.TypeString,
      "progress":    util.TypeInt,
      "weight":      util.TypeInt,
      "missions":    util.TypeAny,
      "exist":       util.TypeBool,
      "editor":      util.TypeString,
      "editTime":    util.TypeTime,
    },
  )

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if name, ok := util.GetString(data, "name"); ok {
    errgo.ErrorIfLenLessThen(name, 4, errgo.ErrProjectNameTooShort)
    errgo.ErrorIfLenMoreThen(name, 15, errgo.ErrProjectNameTooLong)
  }

  // 如果有设置结束时间，将时间改为这天的最晚时间，即23:59:59
  if deadline, ok := util.GetTime(data, "deadline"); ok {
    y1, m1, d1 := deadline.Local().Date()
    t := time.Date(y1, m1, d1, 23, 59, 59, 0, time.Local)
    data["deadline"] = t
    errgo.ErrorIfTimeEarlierThen(t, time.Now(), errgo.ErrProjectDeadlineTooSoon)
  }

  if weight, ok := util.GetInt(data, "weight"); ok {
    if weight != 1 && weight != 2 && weight != 3 {
      return errors.New(errgo.ErrProjectWeightError)
    }
  }

  // 验证每个部门是否正常
  if departments, ok := util.GetAny(data, "departments"); ok {
    var departmentsRef []mgo.DBRef
    departments := departments.([]gjson.Result)
    for _, item := range departments {
      itemId := item.String()
      if item.Exists() && bson.IsObjectIdHex(itemId) {
        _, err := GetDepartmentInfoById(ctx, itemId)
        if err == nil {
          departmentsRef = append(departmentsRef, mgo.DBRef{
            Database:   conf.MgoDBName,
            Collection: model.DepartmentCollection,
            Id:         bson.ObjectIdHex(itemId),
          })
        } else {
          return errors.New(errgo.ErrProjectDepartmentNotFound)
        }
      } else {
        return errors.New(errgo.ErrProjectDepartmentNotFound)
      }
    }

    data["departments"] = departmentsRef
  }

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 先要清缓存，清成功后才可以更新数据
  err := cache.ProjectDel(ctx.Redis, id)

  if err != nil {
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteProjectFailed)
    }
    return errors.New(errgo.ErrUpdateProjectFailed)
  }

  // update
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data["editor.$id"] = bson.ObjectIdHex(ownUserId)
  data["editTime"] = time.Now()

  err = ctx.MgoDB.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrProjectNotFound)
    }
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteProjectFailed)
    }
    return errors.New(errgo.ErrUpdateProjectFailed)
  }

  return nil
}

// 根据id查找项目
func GetProjectInfoById(ctx *context.New, id string) (*model.Project, error) {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrProjectIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Project)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.ProjectGet(ctx.Redis, id, data)

  if !rok {
    err := ctx.MgoDB.C(model.ProjectCollection).Find(bson.M{
      "_id":   bson.ObjectIdHex(id),
      "exist": true,
    }).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrProjectNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.ProjectSet(ctx.Redis, data)
  }

  return data, nil
}

// 查找项目列表(当limit都为0时，查找全部)
func GetProjectsList(ctx *context.New, skip int, limit int, query bson.M) (*model.ProjectList, error) {

  // check
  if skip != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.Project)
  query["exist"] = true

  // find it
  var err error
  if skip == 0 && limit == 0 {
    err = ctx.MgoDB.C(model.ProjectCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = ctx.MgoDB.C(model.ProjectCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrProjectNotFound)
    }
    return nil, err
  }

  // result
  if skip == 0 && limit == 0 {
    return &model.ProjectList{
      List: data,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.ProjectCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrProjectNotFound)
  }

  return &model.ProjectList{
    List:  data,
    Count: count,
    Skip:  skip,
    Limit: limit,
  }, nil
}
