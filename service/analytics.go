package service

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "math"
  "time"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/util"
)

// 获取用户的的基础数据
func GetDepartmentAnalysisById(ctx *context.New, departmentId string) (gin.H, error) {
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
    user     model.User
    missions []model.Mission
  }
  var result = make([]resultStruct, 0)

  // 初始化返回数据的列表
  for _, item := range *userList.List {
    result = append(result, resultStruct{
      user:     item,
      missions: []model.Mission{},
    })
  }

  // 判断用户在该结构体列表内的索引
  var indexOf = func(list []resultStruct, userId bson.ObjectId) int {
    for i, item := range list {
      if item.user.Id == userId {
        return i
      }
    }
    return -1
  }

  // 把任务放到相应的用户内
  for _, item := range *missionList.List {
    if index := indexOf(result, item.User.Id.(bson.ObjectId)); index != -1 {
      result[index].missions = append(result[index].missions, item)
    }
  }

  // 解析返回数据
  var list []gin.H
  for _, item := range result {
    each := item.user.GetMap(model.REMEMBER, "id", "nickname")

    ms := make([]gin.H, 0)
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

  // 返回的数据
  data := department.GetMap(model.REMEMBER, "name", "id")
  data["users"] = list

  return data, nil
}

// 获取部门列表的基础数据
func GetDepartmentsAnalysis(ctx *context.New, skip int, limit int) (gin.H, error) {

  // 返回数据的结构
  type missionStruct struct {
    department model.Department
    missions   []model.Mission
  }

  var result = make([]missionStruct, 0)

  // 获取对应数量的部门
  departmentList, err := GetDepartmentsList(ctx, skip, limit, bson.M{})

  if err != nil {
    return nil, err
  }

  var departmentIdList []bson.ObjectId
  for _, item := range *departmentList.List {
    result = append(result, missionStruct{
      department: item,
      missions:   []model.Mission{},
    })
    departmentIdList = append(departmentIdList, item.Id)
  }

  // 根据这批部门找用户
  userList, err := GetUsersList(ctx, 0, 0, bson.M{
    "department.$id": bson.M{
      "$in": departmentIdList,
    },
  })

  if err != nil {
    return nil, err
  }

  var userIdList []bson.ObjectId
  for _, item := range *userList.List {
    userIdList = append(userIdList, item.Id)
  }

  // 根据这批用户查找任务
  missionList, err := GetMissionsList(ctx, 0, 0, bson.M{
    "user.$id": bson.M{
      "$in": userIdList,
    },
    "status": 1,
  })

  if err != nil {
    return nil, err
  }

  // 根据任务查到用户
  var findUser = func(mission *model.Mission, userList *model.UserList) *model.User {
    for _, item := range *userList.List {
      if mission.User.Id.(bson.ObjectId) == item.Id {
        return &item
      }
    }
    return nil
  }

  // 根据用户查找部门
  var findDepartment = func(user *model.User, departmentList *model.DepartmentList) *model.Department {
    for _, item := range *departmentList.List {
      if user.Department.Id.(bson.ObjectId) == item.Id {
        return &item
      }
    }
    return nil
  }

  // 从返回结果的结构体中找到相应部门的索引
  var indexOf = func(result []missionStruct, departmentId bson.ObjectId) int {
    for i, item := range result {
      if item.department.Id == departmentId {
        return i
      }
    }
    return -1
  }

  for _, item := range *missionList.List {
    if user := findUser(&item, userList); user != nil {
      department := findDepartment(user, departmentList)
      index := indexOf(result, department.Id)
      if result[index].missions == nil {
        result[index].missions = []model.Mission{}
      }
      result[index].missions = append(result[index].missions, item)
    }
  }

  // 返回数据
  list := make([]gin.H, 0)
  for _, item := range result {
    each := item.department.GetMap(model.REMEMBER, "id", "name", "userCount")

    missions := make([]gin.H, 0)
    for _, item := range item.missions {
      each := item.GetMap(model.REMEMBER, "deadline", "id", "name", "progress", "isTimeout")
      missions = append(missions, each)
    }

    each["missions"] = missions

    list = append(list, each)
  }

  return gin.H{
    "list":  list,
    "skip":  skip,
    "limit": limit,
    "count": len(list),
  }, nil
}

// 获取项目列表的基础信息
func GetProjectsAnalysis(ctx *context.New, skip int, limit int) (gin.H, error) {
  projectsList, err := GetProjectsList(ctx, skip, limit, bson.M{})

  if err != nil {
    return nil, err
  }

  // 返回数据
  return projectsList.Each(func(p model.Project) gin.H {
    each := p.GetMap(model.REMEMBER, "isTimeout", "progress", "deadline", "createTime", "name", "id")

    total := p.Deadline.Unix() - p.CreateTime.Unix()
    past := time.Now().Unix() - p.CreateTime.Unix()

    each["totalDay"] = math.Ceil(float64(total) / 60 / 60 / 24)
    each["costDay"] = math.Floor(float64(past) / 60 / 60 / 24)
    each["missionCount"] = len(p.Missions)

    return each
  }), nil
}

// 存用户的数据
func SaveUserAnalysis(m *mgo.Database, day string) error {

  // 找到这天的所有日报
  dailies := new([]model.Daily)

  err := m.C(model.DailyCollection).Find(bson.M{
    "day": day,
  }).All(dailies)

  if err != nil {
    return err
  }

  // 找到所有的人员
  users := new([]model.User)

  err = m.C(model.UserCollection).Find(bson.M{
    "exist": true,
  }).All(users)

  if err != nil {
    return err
  }

  // 找到所有的任务
  missions := new([]model.Mission)

  err = m.C(model.UserCollection).Find(nil).All(missions)

  if err != nil {
    return err
  }

  // 根据用户id找用户
  var findUserById = func(id bson.ObjectId, list *[]model.User) *model.User {
    for _, u := range *list {
      if u.Id == id {
        return &u
      }
    }
    return nil
  }

  // 根据任务id找任务
  var findMissionById = func(id bson.ObjectId, list *[]model.Mission) *model.Mission {
    for _, u := range *list {
      if u.Id == id {
        return &u
      }
    }
    return nil
  }

  // 整理数据
  for _, item := range *dailies {
    user := findUserById(item.User.Id.(bson.ObjectId), users)
    if user != nil {

      dailies := make([]model.UserAnalyticsDaily, 0)

      for _, item := range item.Dailies {
        m := findMissionById(item.MissionId, missions)
        isTimeout := false
        if m != nil {
          isTimeout = m.Deadline.Before(time.Now())
        }
        i := model.UserAnalyticsDaily{
          MissionId:   item.MissionId,
          MissionName: item.MissionName,
          Progress:    item.Progress,
          IsTimeout:   isTimeout,
        }
        dailies = append(dailies, i)
      }

      analysis := model.UserAnalytics{
        User:       item.User,
        Department: user.Department,
        Day:        day,
        Dailies:    dailies,
        CreateTime: time.Now(),
      }

      // 储存出错的话将会一直执行
      util.Forever(func(count int) (done bool) {

        // 超过5次存储，结束
        if count > 5 {
          return true
        }

        c, err := m.C(model.UserAnalyticsCollection).Find(bson.M{
          "day":      day,
          "user.$id": user.Id,
        }).Count()

        if c > 0 {
          m.C(model.UserAnalyticsCollection).Update(bson.M{
            "day":      day,
            "user.$id": user.Id,
          }, analysis)
          // 如果有找着数据，不能重复存
          return true
        } else {
          // 存数据
          err = m.C(model.UserAnalyticsCollection).Insert(analysis)
          // 如果存失败，重来
          if err != nil {
            return false
          }
          return true
        }

        // 其他错误，继续尝试
        return false
      })
    }
  }

  return nil
}
