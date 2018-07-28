package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
  "workerbook/middleware"
)


func registerConsoleDepartmentRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 部门列表
  g.GET("", context.CreateCtx(controller.C_GetDepartmentsList))

  // 获取单个部门
  g.GET("/id/:id", context.CreateCtx(controller.C_GetDepartmentOne))

  // 添加部门
  g.POST("", context.CreateCtx(controller.C_CreateDepartment))

  // 修改部门
  g.PUT("/id/:id", context.CreateCtx(controller.C_UpdateDepartment))

  // 删除部门
  g.DELETE("/id/:id", context.CreateCtx(controller.C_DelDepartmentOne))
}

func registerDepartmentRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 部门列表
  g.GET("", context.CreateCtx(controller.GetDepartmentsList))
}
