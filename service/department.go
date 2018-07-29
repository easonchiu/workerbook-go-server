package service

import (
  "errors"
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

// 创建部门
func CreateDepartment(ctx *context.New, data model.Department) error {

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

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrDepartmentNameEmpty)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 名称不能重复
  count, err := ctx.MgoDB.C(model.DepartmentCollection).Find(bson.M{"name": data.Name}).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateDepartmentFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameDepartmentName)
  }

  // insert it.
  err = ctx.MgoDB.C(model.DepartmentCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateDepartmentFailed)
  }

  return nil
}

// 根据id查找部门信息
func GetDepartmentInfoById(ctx *context.New, id string) (*model.Department, error) {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Department)

  // 先从缓存查数据，找不到的话再从数据库找
  rok := cache.DepartmentGet(ctx.Redis, id, data)
  if !rok {
    err := ctx.MgoDB.C(model.DepartmentCollection).FindId(bson.ObjectIdHex(id)).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrDepartmentNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.DepartmentSet(ctx.Redis, data)
  }

  return data, nil
}

// 查找部门列表(当limit都为0时，查找全部)
func GetDepartmentsList(ctx *context.New, skip int, limit int, query bson.M) (*model.DepartmentList, error) {

  // check
  if limit != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.Department)
  query["exist"] = true

  // find it
  var err error
  if skip == 0 && limit == 0 {
    err = ctx.MgoDB.C(model.DepartmentCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = ctx.MgoDB.C(model.DepartmentCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrDepartmentNotFound)
    }
    return nil, err
  }

  // result

  if skip == 0 && limit == 0 {
    return &model.DepartmentList{
      List: data,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.DepartmentCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrDepartmentNotFound)
  }

  return &model.DepartmentList{
    List:  data,
    Count: count,
    Skip:  skip,
    Limit: limit,
  }, nil
}

// 全量更新所有部门的人数
func UpdateDepartmentsUserCount(ctx *context.New) error {

  departments := new([]model.Department)
  err := ctx.MgoDB.C(model.DepartmentCollection).Find(bson.M{"exist": true}).All(departments)

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  for _, department := range *departments {
    count, err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
      "exist":          true,
      "department.$id": department.Id,
    }).Count()

    if err != nil {
      return err
    }

    // 清空缓存并更新数据
    if err := cache.DepartmentDel(ctx.Redis, department.Id.Hex()); err == nil {
      ctx.MgoDB.C(model.DepartmentCollection).UpdateId(department.Id, bson.M{
        "$set": bson.M{
          "userCount": count,
        },
      })
    }
  }

  return nil
}

// 更新部门信息
func UpdateDepartment(ctx *context.New, id string, data bson.M) error {

  if data == nil {
    return errors.New(errgo.ErrServerError)
  }

  // 限制更新字段
  util.Only(
    data,
    util.Keys{
      "name":      util.TypeString,
      "userCount": util.TypeInt,
      "exist":     util.TypeBool,
      "editor":    util.TypeBsonM,
      "editTime":  util.TypeTime,
    },
  )

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 名称唯一
  if name, ok := util.GetString(data, "name"); ok {
    count, err := ctx.MgoDB.C(model.DepartmentCollection).Find(bson.M{
      "name":  name,
      "exist": true,
      "_id": bson.M{
        "$ne": bson.ObjectIdHex(id),
      },
    }).Count()

    if err != nil {
      return errors.New(errgo.ErrUpdateDepartmentFailed)
    }

    if count > 0 {
      return errors.New(errgo.ErrSameDepartmentName)
    }
  }

  // 更新前必须成功清空缓存
  err := cache.DepartmentDel(ctx.Redis, id)

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  // 更新数据
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data["editor.$id"] = bson.ObjectIdHex(ownUserId)
  data["editTime"] = time.Now()

  err = ctx.MgoDB.C(model.DepartmentCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrDepartmentNotFound)
    }
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  return nil
}

// 根据id删除部门（前提是部门内无人）
func DelDepartmentById(ctx *context.New, id string) error {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 查找该部门内是否有用户
  count, err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
    "exist":          true,
    "department.$id": bson.ObjectIdHex(id),
  }).Count()

  if count > 0 {
    return errors.New(errgo.ErrDeleteDepartmentHasUser)
  }

  // 清除缓存，缓存清成功才可以清数据，不然会有脏数据
  err = cache.DepartmentDel(ctx.Redis, id)

  if err != nil {
    return errors.New(errgo.ErrDeleteDepartmentFailed)
  }

  // 拿操作人id
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  // 删除项目里相关的部门，和缓存
  projects := new([]model.Project)

  err = ctx.MgoDB.C(model.ProjectCollection).Find(bson.M{
    "departments.$id": bson.ObjectIdHex(id),
  }).All(projects)

  if err != nil {
    return errors.New(errgo.ErrDeleteDepartmentFailed)
  }

  for _, item := range *projects {
    cache.ProjectDel(ctx.Redis, item.Id.Hex())
  }

  if len(*projects) > 0 {
    err = ctx.MgoDB.C(model.ProjectCollection).Update(bson.M{
      "departments.$id": bson.ObjectIdHex(id),
    }, bson.M{
      "$pull": bson.M{
        "departments": bson.M{
          "$id": bson.ObjectIdHex(id),
        },
      },
    })

    if err != nil {
      return errors.New(errgo.ErrDeleteDepartmentFailed)
    }
  }

  // 删除部门
  err = ctx.MgoDB.C(model.DepartmentCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": bson.M{
      "exist":      false,
      "editor.$id": bson.ObjectIdHex(ownUserId),
      "editTime":   time.Now(),
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrDepartmentNotFound)
    }
    return errors.New(errgo.ErrDeleteDepartmentFailed)
  }

  return nil
}
