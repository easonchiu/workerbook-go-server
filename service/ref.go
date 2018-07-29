package service

import (
  "errors"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "workerbook/cache"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
)

// 根据mgo的ref查找数据
// 如果缓存有，则从缓存获取
// 如果缓存没有，从数据库获取，且正常获取后存到缓存

func FindDepartmentRef(ctx *context.New, ref *mgo.DBRef) (*model.Department, error) {
  if ref.Id == nil {
    return nil, errors.New(errgo.ErrDepartmentIdError)
  }

  var department = new(model.Department)

  if ok := cache.DepartmentGet(ctx.Redis, ref.Id.(bson.ObjectId).Hex(), department); !ok {
    err := ctx.MgoDB.C(ref.Collection).Find(bson.M{
      "_id":   ref.Id.(bson.ObjectId),
      "exist": true,
    }).One(department)

    if err == nil {
      cache.DepartmentSet(ctx.Redis, department)
      return department, nil
    }

    return department, errors.New(errgo.ErrDepartmentNotFound)
  }

  return department, nil
}

func FindUserRef(ctx *context.New, ref *mgo.DBRef) (*model.User, error) {
  if ref.Id == nil {
    return nil, errors.New(errgo.ErrUserIdError)
  }

  var user = new(model.User)

  if ok := cache.UserGet(ctx.Redis, ref.Id.(bson.ObjectId).Hex(), user); !ok {
    err := ctx.MgoDB.C(ref.Collection).Find(bson.M{
      "_id": ref.Id.(bson.ObjectId),
      // "exist": true,
    }).One(user)

    if err == nil {
      cache.UserSet(ctx.Redis, user)
      return user, nil
    }

    return user, errors.New(errgo.ErrUserNotFound)
  }

  return user, nil
}

func FindProjectRef(ctx *context.New, ref *mgo.DBRef) (*model.Project, error) {
  if ref.Id == nil {
    return nil, errors.New(errgo.ErrProjectIdError)
  }

  var project = new(model.Project)

  if ok := cache.ProjectGet(ctx.Redis, ref.Id.(bson.ObjectId).Hex(), project); !ok {
    err := ctx.MgoDB.C(ref.Collection).Find(bson.M{
      "_id":   ref.Id.(bson.ObjectId),
      "exist": true,
    }).One(project)

    if err == nil {
      cache.ProjectSet(ctx.Redis, project)
      return project, nil
    }

    return project, errors.New(errgo.ErrProjectNotFound)
  }

  return project, nil
}

func FindMissionRef(ctx *context.New, ref *mgo.DBRef) (*model.Mission, error) {
  if ref.Id == nil {
    return nil, errors.New(errgo.ErrMissionIdError)
  }

  var mission = new(model.Mission)

  if ok := cache.MissionGet(ctx.Redis, ref.Id.(bson.ObjectId).Hex(), mission); !ok {
    err := ctx.MgoDB.C(ref.Collection).Find(bson.M{
      "_id":   ref.Id.(bson.ObjectId),
      "exist": true,
    }).One(mission)

    if err == nil {
      cache.MissionSet(ctx.Redis, mission)
      return mission, nil
    }

    return mission, errors.New(errgo.ErrMissionNotFound)
  }

  return mission, nil
}
