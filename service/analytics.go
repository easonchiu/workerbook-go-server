package service

import (
  "github.com/easonchiu/repego"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
)

// 获取部门列表的基础数据
func GetDepartmentsListAnalysis(ctx *context.New, skip int, limit int) (*model.DepartmentAnalyticsList, error) {

  // 初始化返回值结构体
  var result = new(model.DepartmentAnalyticsList)

  // 获取对应数量的部门
  departmentList, err := GetDepartmentsList(ctx, skip, limit, bson.M{})

  if err != nil {
    return nil, err
  }

  // 根据这批部门找用户
  userList, err := GetUsersList(ctx, 0, 0, bson.M{
    "department.$id": bson.M{
      "$in": departmentList.Ids(),
    },
    "exist": true,
  })

  if err != nil {
    return nil, err
  }

  // 根据这批用户查找任务
  missionList, err := GetMissionsList(ctx, 0, 0, bson.M{
    "user.$id": bson.M{
      "$in": userList.Ids(),
    },
    "status": 1,
  })

  if err != nil {
    return nil, err
  }

  // 初始化返回结果
  departmentList.Each(func(department *model.Department) gin.H {
    result.List = append(result.List, &model.DepartmentAnalytics{
      Department: department,
      Missions:   []*model.Mission{},
    })
    return nil
  })

  // 根据任务查到用户，然后再把任务数据追加到相应用户的统计数据中
  missionList.Each(func(mission *model.Mission) gin.H {
    // 从用户列表中找相应的用户
    if user := userList.Find(mission.User.Id.(bson.ObjectId)); user != nil {
      // 从部门列表中找到相应的部门
      if department := departmentList.Find(user.Department.Id.(bson.ObjectId)); department != nil {
        // 找到部门并追加任务数据
        if d := result.Find(department.Id); d != nil {
          if d.Missions == nil {
            d.Missions = []*model.Mission{}
          }

          d.Missions = append(d.Missions, mission)
        }
      }
    }

    return nil
  })

  result.Skip = skip
  result.Limit = limit
  result.Count = departmentList.Count

  return result, nil
}

// 获取项目列表的基础信息
func GetProjectsListAnalysis(ctx *context.New, skip int, limit int) (*model.ProjectListAnalytics, error) {

  // 初始化返回结果的结构体
  var result = new(model.ProjectListAnalytics)

  // 找到相应数量的项目
  projectsList, err := GetProjectsList(ctx, skip, limit, bson.M{})

  if err != nil {
    return nil, err
  }

  // 返回数据
  projectsList.Each(func(p *model.Project) gin.H {

    item := model.ProjectAnalytics{
      Project:  p,
      Missions: make([]*model.Mission, 0, len(p.Missions)),
    }

    for _, m := range p.Missions {
      mission, err := FindMissionRef(ctx, &m)
      if err == nil {
        item.Missions = append(item.Missions, mission)
      }
    }

    result.List = append(result.List, &item)

    return nil
  })

  result.Skip = skip
  result.Limit = limit
  result.Count = projectsList.Count

  return result, nil
}

// 获取部门的概要数据
func GetDepartmentSummaryAnalysisById(ctx *context.New, departmentId string) (*model.DepartmentUsersAnalytics, error) {
  // catch
  ctx.Errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)

  if err := ctx.Errgo.PopError(); err != nil {
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

  // 获取所有和该部门用户相关的正常任务
  missionList, err := GetMissionsList(ctx, 0, 0, bson.M{
    "user.$id": bson.M{
      "$in": userList.Ids(),
    },
    "status": 1,
  })

  if err != nil {
    return nil, err
  }

  // 初始化返回数据的列表
  var result = new(model.DepartmentUsersAnalytics)

  result.Department = department

  userList.Each(func(user *model.User) gin.H {
    result.Users = append(result.Users, &model.UserAnalytics{
      User:     user,
      Missions: []*model.Mission{},
    })

    return nil
  })

  // 把任务放到相应的用户内
  missionList.Each(func(mission *model.Mission) gin.H {
    if u := result.Find(mission.User.Id.(bson.ObjectId)); u != nil {
      u.Missions = append(u.Missions, mission)
    }

    return nil
  })

  return result, nil
}

// 获取部门的详细数据
func GetDepartmentDetailAnalysisById(ctx *context.New, departmentId string) ([]*model.UserMissionsChart, error) {
  // catch
  ctx.Errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  // 获取所有部门成员
  users, err := GetUsersList(ctx, 0, 0, bson.M{
    "department.$id": bson.ObjectIdHex(departmentId),
  })

  if err != nil {
    return nil, err
  }

  // 找到部门用户的所有任务
  missions, err := GetMissionsList(ctx, 0, 0, bson.M{
    "user.$id": bson.M{
      "$in": users.Ids(),
    },
  })

  if err != nil {
    return nil, err
  }

  projects, err := GetProjectsList(ctx, 0, 0, bson.M{})

  if err != nil {
    return nil, err
  }

  // 查找所有任务相关的数据
  analytics := new([]model.Analytics)
  err = ctx.MgoDB.C(model.AnalyticsCollection).Find(bson.M{
    "missionId": bson.M{
      "$in": missions.Ids(),
    },
  }).All(analytics)

  if err != nil {
    return nil, err
  }

  // 初始化返回数据
  var result = make([]*model.UserMissionsChart, 0)

  for _, item := range users.List {

    // 初始化任务
    mission := make([]*model.MissionChartData, 0)

    for _, m := range missions.List {
      if m.User.Id.(bson.ObjectId) == item.Id {

        // 找到这个任务的项目
        project := projects.Find(m.Project.Id.(bson.ObjectId))

        if project != nil {
          var mcd = &model.MissionChartData{
            Id:          m.Id,
            Name:        m.Name,
            ProjectId:   project.Id,
            ProjectName: project.Name,
            IsTimeout:   m.Deadline.Before(time.Now()),
            Deadline:    m.Deadline,
            Data:        make([]*model.MissionChartDay, 0, 120),
          }

          // 遍历所有统计数据，把对应任务的放进去
          for _, a := range *analytics {
            if a.MissionId == m.Id {
              mcd.Append(&model.MissionChartDay{
                Progress:  a.Progress,
                ChartTime: a.CreateTime,
                Day:       a.Day,
              })
            }
          }

          mission = append(mission, mcd)
        }
      }
    }

    result = append(result, &model.UserMissionsChart{
      User:   item,
      Charts: mission,
    })
  }

  return result, nil
}

// 获取项目的概要数据
func GetProjectSummaryAnalysisById(ctx *context.New, projectId string) (*model.ProjectMissionsChart, error) {
  // catch
  ctx.Errgo.ErrorIfStringNotObjectId(projectId, errgo.ErrDepartmentIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  // 找到项目数据
  project, err := GetProjectInfoById(ctx, projectId)

  if err != nil {
    return nil, err
  }

  // 找到这个项目的所有任务
  missions, err := GetMissionsList(ctx, 0, 0, bson.M{
    "project.$id": project.Id,
  })

  if err != nil {
    return nil, err
  }

  // 找到数据
  analytics := new([]model.Analytics)
  err = ctx.MgoDB.C(model.AnalyticsCollection).Find(bson.M{
    "projectId": project.Id,
  }).Sort("-createTime").All(analytics)

  if err != nil {
    return nil, err
  }

  // 初始化返回数据的结构
  var result = new(model.ProjectMissionsChart)

  result.Project = project

  for _, item := range missions.List {
    result.Charts = append(result.Charts, &model.MissionChartData{
      Id:          item.Id,
      Name:        item.Name,
      Deadline:    item.Deadline,
      ProjectId:   project.Id,
      ProjectName: project.Name,
      IsTimeout:   item.Deadline.Before(time.Now()),
      Data:        []*model.MissionChartDay{},
    })
  }

  // 遍历查到的数据
  for _, item := range *analytics {
    // 找到任务
    res := result.Find(item.MissionId)

    if res != nil {
      // 添加数据
      res.Append(&model.MissionChartDay{
        Progress:  item.Progress,
        ChartTime: item.CreateTime,
        Day:       item.Day,
      })
    }
  }

  return result, nil
}

// 获取项目的详细数据
func GetProjectDetailAnalysisById(ctx *context.New, projectId string) (*model.ProjectMissionsChart, error) {
  // catch
  ctx.Errgo.ErrorIfStringNotObjectId(projectId, errgo.ErrDepartmentIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  // 找到该项目
  project, err := GetProjectInfoById(ctx, projectId)

  if err != nil {
    return nil, err
  }

  // 找到该项目的所有任务
  missions, err := GetMissionsList(ctx, 0, 0, bson.M{
    "project.$id": bson.ObjectIdHex(projectId),
  })

  if err != nil {
    return nil, err
  }

  // 查找所有任务相关的数据
  analytics := new([]model.Analytics)

  err = ctx.MgoDB.C(model.AnalyticsCollection).Find(bson.M{
    "projectId": bson.ObjectIdHex(projectId),
  }).All(analytics)

  if err != nil {
    return nil, err
  }

  // 初始化返回数据
  var result = new(model.ProjectMissionsChart)

  // 初始化任务
  mission := make([]*model.MissionChartData, 0)

  for _, m := range missions.List {
    var mcd = &model.MissionChartData{
      Id:          m.Id,
      Name:        m.Name,
      ProjectId:   project.Id,
      ProjectName: project.Name,
      IsTimeout:   m.Deadline.Before(time.Now()),
      Deadline:    m.Deadline,
      Data:        make([]*model.MissionChartDay, 0, 120),
    }

    // 遍历所有统计数据，把对应任务的放进去
    for _, a := range *analytics {
      if a.MissionId == m.Id {
        mcd.Append(&model.MissionChartDay{
          Progress:  a.Progress,
          ChartTime: a.CreateTime,
          Day:       a.Day,
        })
      }
    }

    mission = append(mission, mcd)
  }

  result.Project = project
  result.Charts = mission

  return result, nil
}

// 存数据
func SaveAnalysisByDay(m *mgo.Database, day time.Time) error {

  // 找到所有未结束的项目
  var projects model.ProjectList

  err := m.C(model.ProjectCollection).Find(bson.M{
    "exist":  true,
    "status": 1,
  }).All(&projects.List)

  if err != nil {
    return err
  }

  // 根据这些项目筛选出相应的任务(进度为0的不存)
  var missions model.MissionList

  err = m.C(model.MissionCollection).Find(bson.M{
    "project.$id": bson.M{
      "$in": projects.Ids(),
    },
    "progress": bson.M{
      "$ne": 0,
    },
    "chartTime": bson.M{
      "$ne": day.Format("2006-01-02"),
    },
    "exist": true,
  }).All(&missions.List)

  if err != nil {
    return err
  }

  // 找到所有的人员
  var users model.UserList

  err = m.C(model.UserCollection).Find(bson.M{
    "exist": true,
  }).All(&users.List)

  if err != nil {
    return err
  }

  // 存数据
  missions.Each(func(mission *model.Mission) gin.H {

    // 找到用户数据
    user := users.Find(mission.User.Id.(bson.ObjectId))

    if user != nil {
      analytics := model.Analytics{
        MissionId:  mission.Id,
        Progress:   mission.Progress,
        ProjectId:  mission.Project.Id.(bson.ObjectId),
        Day:        day.Format("2006-01-02"),
        CreateTime: time.Now(),
      }

      // 储存出错的话将会连续执行5次，5次过后，写错误日志
      repego.Call(func(r *repego.R) {

        // 超过5次存储，结束
        if r.Count > 5 {
          r.Done()
          return
        }

        // 存数据
        err := m.C(model.AnalyticsCollection).Insert(analytics)

        // ok
        if err == nil {

          // 这里要想办法确保更新不出问题
          repego.Call(func(r *repego.R) {
            // 超过5次存储，结束
            if r.Count > 5 {
              r.Done()
              return
            }

            err := m.C(model.MissionCollection).UpdateId(analytics.MissionId, bson.M{
              "$set": bson.M{
                "chartTime":   day,
                "preProgress": mission.Progress,
              },
            })

            // ok
            if err == nil {
              r.Done()
            }
          }).Do(time.Second) // 间隔时间为1s

          r.Done()
        }

      }).Do(time.Second) // 间隔时间为1s

    }
    return nil
  })

  return nil
}
