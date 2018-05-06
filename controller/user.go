package controller

import (
	"github.com/gin-gonic/gin"
	"web/service"
)

func GetUserInfo(c *gin.Context) {
	resp := Response{c}

	id := c.Params.ByName("id")

	userInfo, err := service.GetUserInfoById(id)
	if err != nil {
		resp.Error(err)
		return
	}

	resp.Success(userInfo)
}