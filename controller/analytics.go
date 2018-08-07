package controller

import (
  "github.com/gin-gonic/gin"
  "math"
  "time"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个部门成员的概要信息
func GetDepartmentSummaryAnalytics(ctx *context.New) {
  // get
  departmentId, _ := ctx.GetParam("id")

  data, err := service.GetDepartmentSummaryAnalysisById(ctx, departmentId)

  if err != nil {
    ctx.Error(err)
    return
  }

  // 解析返回数据
  result := data.Department.GetMap(model.REMEMBER, "name", "id")

  var users []gin.H
  for _, item := range data.Users {
    each := item.User.GetMap(model.REMEMBER, "id", "nickname")

    missions := make([]gin.H, 0, len(item.Missions))
    for _, m := range item.Missions {
      data := m.GetMap(model.REMEMBER, "deadline", "id", "name", "progress", "isTimeout")
      missions = append(missions, data)
    }

    each["missions"] = missions

    users = append(users, each)
  }

  result["users"] = users

  ctx.Success(gin.H{
    "data": result,
  })
}

// 获取部门列表的统计信息
func GetDepartmentListAnalytics(ctx *context.New) {
  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetDepartmentsListAnalysis(ctx, skip, limit)

  if err != nil {
    ctx.Error(err)
    return
  }

  // 返回数据
  list := make([]gin.H, 0, data.Count)

  for _, item := range data.List {
    each := item.Department.GetMap(model.REMEMBER, "id", "name", "userCount")

    missions := make([]gin.H, 0, len(item.Missions))
    for _, item := range item.Missions {
      each := item.GetMap(model.REMEMBER, "deadline", "id", "name", "progress", "isTimeout")
      missions = append(missions, each)
    }

    each["missions"] = missions

    list = append(list, each)
  }

  if skip == 0 && limit == 0 {
    ctx.Success(gin.H{
      "data": gin.H{
        "list": list,
      },
    })
    return
  }

  ctx.Success(gin.H{
    "data": gin.H{
      "list":  list,
      "skip":  skip,
      "limit": limit,
      "count": data.Count,
    },
  })
}

// 获取部门的详细数据
func GetDepartmentDetailAnalytics(ctx *context.New) {
  // get
  departmentId, _ := ctx.GetParam("id")

  data, err := service.GetDepartmentDetailAnalysisById(ctx, departmentId)

  if err != nil {
    ctx.Error(err)
  }

  result := make([]gin.H, 0)

  for _, item := range data {
    each := item.User.GetMap(model.REMEMBER, "id", "nickname", "title")

    charts := make([]gin.H, 0, len(item.Charts))
    for _, i := range item.Charts {
      charts = append(charts, i.GetMap())
    }

    each["missions"] = charts

    result = append(result, each)
  }

  ctx.Success(gin.H{
    "data": result,
  })
}

// 获取项目列表的统计信息
func GetProjectListAnalytics(ctx *context.New) {
  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetProjectsListAnalysis(ctx, skip, limit)

  if err != nil {
    ctx.Error(err)
    return
  }

  list := make([]gin.H, 0, data.Count)

  for _, item := range data.List {
    each := item.Project.GetMap(model.REMEMBER, "isTimeout", "progress", "deadline", "createTime", "name", "id")
    each["progress"] = service.GetProjectProgress(ctx, item.Project.Id)

    missions := make([]gin.H, 0, len(item.Missions))
    for _, item := range item.Missions {
      each := item.GetMap(model.REMEMBER, "id", "isTimeout")
      missions = append(missions, each)
    }

    each["missions"] = missions

    total := item.Project.Deadline.Unix() - item.Project.CreateTime.Unix()
    past := time.Now().Unix() - item.Project.CreateTime.Unix()

    each["totalDay"] = math.Ceil(float64(total) / 60 / 60 / 24)
    each["costDay"] = math.Floor(float64(past) / 60 / 60 / 24)
    each["missionCount"] = len(item.Missions)

    list = append(list, each)
  }

  if skip == 0 && limit == 0 {
    ctx.Success(gin.H{
      "data": gin.H{
        "list": list,
      },
    })
    return
  }

  ctx.Success(gin.H{
    "data": gin.H{
      "list":  list,
      "skip":  skip,
      "limit": limit,
      "count": data.Count,
    },
  })
}

// 获取单个项目的任务概要信息
func GetProjectSummaryAnalytics(ctx *context.New) {
  // get
  projectId, _ := ctx.GetParam("id")

  data, err := service.GetProjectSummaryAnalysisById(ctx, projectId)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data.GetMap(),
  })
}

// 获取项目的详细数据
func GetProjectDetailAnalytics(ctx *context.New) {
  // get
  projectId, _ := ctx.GetParam("id")

  data, err := service.GetProjectDetailAnalysisById(ctx, projectId)

  if err != nil {
    ctx.Error(err)
  }

  result := data.Project.GetMap(model.REMEMBER, "id", "name", "deadline")

  charts := make([]gin.H, 0, len(data.Charts))
  for _, i := range data.Charts {
    charts = append(charts, i.GetMap())
  }

  result["missions"] = charts

  ctx.Success(gin.H{
    "data": result,
  })
}