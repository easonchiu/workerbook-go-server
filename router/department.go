package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerDepartmentRouter(g *gin.RouterGroup) {

  // 部门列表
  g.GET("", middleware.ConsoleJwt, controller.GetDepartmentsList)

  // 全部部门列表
  g.GET("all", middleware.ConsoleJwt, controller.GetAllDepartmentsList)

  // 添加部门
  g.POST("", middleware.ConsoleJwt, controller.CreateDepartment)

  // 修改部门
  g.PUT("/:id", middleware.ConsoleJwt, controller.UpdateDepartment)

}
