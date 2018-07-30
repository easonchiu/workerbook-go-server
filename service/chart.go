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

  // 获取所有正常任务
  missionList, err := GetMissionsList(ctx, 0, 20, bson.M{
    // "user.$id": bson.M{
    //   "$in": userIdList,
    // },
    "status": 1,
  })

  if err != nil {
    return nil, err
  }

  // 返回结果的数据结构
  type resultType map[mgo.DBRef][]model.Mission
  var result = make(resultType)

  for _, item := range *userList.List {
    ref := mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         item.Id,
    }
    result[ref] = []model.Mission{}
  }

  for _, item := range *missionList.List {
    if _, ok := result[item.User]; ok {
      result[item.User] = append(result[item.User], item)
    }
  }

  // 解析返回数据
  var list []gin.H
  for user, missions := range result {
    u, err := FindUserRef(ctx, &user)

    if err != nil {
      break
    }

    item := u.GetMap(model.REMEMBER, "id", "nickname")

    ms := []gin.H{}
    for _, m := range missions {
      data := m.GetMap(model.REMEMBER, "deadline", "id", "name", "project", "progress", "isTimeout")
      project, err := FindProjectRef(ctx, &m.Project)
      if err == nil {
        data["project"] = project.GetMap(model.REMEMBER, "name", "weight")
      }
      ms = append(ms, data)
    }
    item["missions"] = ms

    list = append(list, item)
  }

  return gin.H{
    "department": department.GetMap(model.REMEMBER, "name", "id"),
    "list": list,
  }, nil
}
