package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleRouter(g  *gin.RouterGroup) {

  // 用户列表
  g.GET("/users", middleware.ConsoleJwt, controller.GetUsersList)

  // 添加用户
  g.POST("/users", middleware.ConsoleJwt, controller.CreateUser)

  // 修改用户
  g.PUT("/users/:id", middleware.ConsoleJwt, controller.UpdateUser)


  // 部门列表
  g.GET("/departments", middleware.ConsoleJwt, controller.GetDepartmentsList)

  // 添加部门
  g.POST("/departments", middleware.ConsoleJwt, controller.CreateDepartment)

  // 修改部门
  g.PUT("/departments/:id", middleware.ConsoleJwt, controller.UpdateDepartment)

}
