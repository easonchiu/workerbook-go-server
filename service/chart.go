package service

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
)

func GetDepartmentUserSummary(ctx *context.New, departmentId string) (gin.H, error) {
  // catch
  errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  // 获取部门信息
  department, err := GetDepartmentInfoById(ctx, departmentId)

  if err != nil {
    return nil, err
  }

  // 获取所有部门成员
  userList, err := GetUsersList(ctx, 0, 0, bson.M{
    "department.$id": bson.ObjectIdHex(departmentId),
  })

  if err != nil {
    return nil, err
  }

  var userIdList []bson.ObjectId
  for _, item := range *userList.List {
    userIdList = append(userIdList, item.Id)
  }

  // 获取所有和该部门用户相关的正常任务
  missionList, err := GetMissionsList(ctx, 0, 20, bson.M{
    "user.$id": bson.M{
      "$in": userIdList,
    },
    "status": 1,
  })

  if err != nil {
    return nil, err
  }

  // 返回结果的数据结构
  type resultStruct struct {
    user     mgo.DBRef
    missions []model.Mission
  }
  var result = make([]resultStruct, 0)

  // 初始化返回数据的列表
  for _, item := range *userList.List {
    ref := mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         item.Id,
    }
    result = append(result, resultStruct{
      user:     ref,
      missions: []model.Mission{},
    })
  }

  // 判断用户在该结构体列表内的索引
  var indexOf = func(list []resultStruct, user mgo.DBRef) int {
    i := -1
    for index, item := range list {
      if item.user == user {
        return index
      }
    }
    return i
  }

  // 把任务放到相应的用户内
  for _, item := range *missionList.List {
    if index := indexOf(result, item.User); index != -1 {
      result[index].missions = append(result[index].missions, item)
    }
  }

  // 解析返回数据
  var list []gin.H
  for _, item := range result {
    u, err := FindUserRef(ctx, &item.user)

    if err != nil {
      break
    }

    each := u.GetMap(model.REMEMBER, "id", "nickname")

    ms := []gin.H{}
    for _, m := range item.missions {
      data := m.GetMap(model.REMEMBER, "deadline", "id", "name", "project", "progress", "isTimeout")
      project, err := FindProjectRef(ctx, &m.Project)
      if err == nil {
        data["project"] = project.GetMap(model.REMEMBER, "name", "weight")
      }
      ms = append(ms, data)
    }
    each["missions"] = ms

    list = append(list, each)
  }

  return gin.H{
    "department": department.GetMap(model.REMEMBER, "name", "id"),
    "list":       list,
  }, nil
}
