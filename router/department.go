package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerDepartmentRouter(g *gin.RouterGroup) {

  // 部门列表
  g.GET("", middleware.ConsoleJwt, controller.GetDepartmentsList)

  // 获取单个部门
  g.GET("/:id", middleware.Jwt, controller.GetDepartmentOne)

  // 添加部门
  g.POST("", middleware.ConsoleJwt, controller.CreateDepartment)

  // 修改部门
  g.PUT("/:id", middleware.ConsoleJwt, controller.UpdateDepartment)

}
