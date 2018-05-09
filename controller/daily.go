package controller

import (
	"github.com/gin-gonic/gin"
	"workerbook/service"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"workerbook/model"
)

// 获取日报列表
func GetDailiesList(c *gin.Context) {
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

	dailiesList, err := service.GetDailiesList(intSkip, intLimit)
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(dailiesList)
}

// 获取单个日报的信息
func GetDailyInfo(c *gin.Context) {
	ctx := CreateCtx(c)

	id := ctx.getParam("id")

	dailyInfo, err := service.GetDailyInfoById(bson.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(dailyInfo)
}

// 创建日报
func CreateDaily(c *gin.Context) {
	ctx := CreateCtx(c)

	err := service.CreateDaily(model.Daily{})
	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(nil)
}