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

	ctx.Success(gin.H{
		"list": dailiesList,
	})
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

	ctx.Success(gin.H{
		"data": dailyInfo,
	})
}

// 创建日报
func CreateDailyItem(c *gin.Context) {
	ctx := CreateCtx(c)

	id := ctx.getRaw("id")

	// 找到用户今天的日报内容(找不到会创建一个空内容的日报数据)
	dailyInfo, err := service.GetUserTodayDailyByUid(bson.ObjectIdHex(id))

	if err != nil {
		ctx.Error(err, 1)
		return
	}

	// 一条日报数据
	data := model.DailyItem{
		Id: bson.NewObjectId(),
		Record: "写了啥写了啥",
		Progress: 50,
		Pname: "某项目",
		Pid: "5af501c4421aa996bd7a7733",
	}

	// 插入数据
	err = service.AppendDailyItemIntoUsersDailyList(data, dailyInfo.Id)

	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(gin.H{
		"data": dailyInfo,
	})
}

// 删除日报
func DeleteDailyItem(c *gin.Context) {
	ctx := CreateCtx(c)

	ctx.Success(gin.H{})
}

// 获取今天的日报
func GetTodayDaily(c *gin.Context) {
	ctx := CreateCtx(c)

	uid := ctx.getParam("uid")

	dailyInfo, err := service.GetUserTodayDailyByUid(bson.ObjectIdHex(uid))

	if err != nil {
		ctx.Error(err, 1)
		return
	}

	ctx.Success(gin.H{
		"data": dailyInfo,
	})
}