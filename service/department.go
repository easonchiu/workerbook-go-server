package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
)

// 创建部门
func CreateDepartment(ctx *context.Context, data model.Department) error {

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

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
func GetDepartmentInfoById(ctx *context.Context, id string) (gin.H, error) {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Department)

  // 先从缓存查数据，找不到的话再从数据库找
  rok := cache.DepartmentGet(id, data)
  if !rok {
    err := ctx.MgoDB.C(model.DepartmentCollection).FindId(bson.ObjectIdHex(id)).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrDepartmentNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.DepartmentSet(id, data)
  }

  return data.GetMap(FindRef(ctx.MgoDB)), nil
}

// 查找部门列表(当skip和limit都为0时，查找全部)
func GetDepartmentsList(ctx *context.Context, skip int, limit int, query bson.M) (gin.H, error) {

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
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(FindRef(ctx.MgoDB)))
  }

  if skip == 0 && limit == 0 {
    return gin.H{
      "list": list,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.DepartmentCollection).Find(query).Count()

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
func UpdateDepartmentsUserCount(ctx *context.Context) error {

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
    if err := cache.DepartmentDel(department.Id.Hex()); err != nil {
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
func UpdateDepartment(ctx *context.Context, id string, data bson.M) error {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 名称唯一
  if name, ok := data["name"]; ok {
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
  err := cache.DepartmentDel(id)

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  // 更新数据
  err = ctx.MgoDB.C(model.DepartmentCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateDepartmentFailed)
  }

  return nil
}

// 根据id删除部门（前提是部门内无人）
func DelDepartmentById(ctx *context.Context, id string) error {

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
  err = cache.DepartmentDel(id)

  if err != nil {
    return errors.New(errgo.ErrDeleteDepartmentFailed)
  }

  // 删除
  err = ctx.MgoDB.C(model.DepartmentCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": bson.M{
      "exist": false,
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrDepartmentNotFound)
    }
    return err
  }

  if err != nil {
    return errors.New(errgo.ErrDeleteDepartmentFailed)
  }

  return nil
}
