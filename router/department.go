package router

import (
  `github.com/gin-gonic/gin`
  "workerbook/controller"
  "workerbook/middleware"
)


func registerConsoleDepartmentRouter(g *gin.RouterGroup) {
  // 部门列表
  g.GET("", middleware.ConsoleJwt, controller.C_GetDepartmentsList)

  // 获取单个部门
  g.GET("/:id", middleware.ConsoleJwt, controller.C_GetDepartmentOne)

  // 添加部门
  g.POST("", middleware.ConsoleJwt, controller.C_CreateDepartment)

  // 修改部门
  g.PUT("/:id", middleware.ConsoleJwt, controller.C_UpdateDepartment)

  // 删除部门
  g.DELETE("/:id", middleware.ConsoleJwt, controller.C_DelDepartmentOne)
}

func registerDepartmentRouter(g *gin.RouterGroup) {

}
