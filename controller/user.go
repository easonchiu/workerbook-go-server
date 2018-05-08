package controller

import (
	"github.com/gin-gonic/gin"
	"workerbook/service"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"workerbook/model"
	"time"
	"fmt"
)

// 获取用户列表
func GetUsersList(c *gin.Context) {
	ctx := CreateCtx(c)

	skip, _ := c.GetQuery("skip")
	limit, _  := c.GetQuery("limit")

	intSkip, err := strconv.Atoi(skip)

	if err != nil {
		intSkip = 0
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		intLimit = 10
	}

	userInfo, err := service.GetUsersList(intSkip, intLimit)
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(userInfo)
}

// 获取单个用户的信息
func GetUserInfo(c *gin.Context) {
	ctx := CreateCtx(c)

	id := ctx.getParam("id")

	userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(userInfo)
}

// 创建用户
func CreateUser(c *gin.Context) {
	ctx := CreateCtx(c)

	data := model.User{
		NickName: ctx.getRaw("nickname"),
		Email: ctx.getRaw("email"),
		UserName: ctx.getRaw("username"),
		Gid: ctx.getRaw("gid"),
		Mobile: ctx.getRaw("mobile"),
		Password: ctx.getRaw("password"),
		Role: 1,
		CreateTime: time.Now(),
	}

	fmt.Println(data)

	err := service.CreateUser(data)
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(nil)
}