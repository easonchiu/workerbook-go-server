package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tidwall/gjson"
)

type Context struct {
	Ctx *gin.Context
	RawData []byte
}

// create a new context
func CreateCtx(c *gin.Context) *Context {
	bytes, _ := c.GetRawData()

	ctx := &Context{
		Ctx: c,
		RawData: bytes,
	}

	return ctx
}

// success handle
func (c *Context) Success(data gin.H) {
	respH := gin.H{
		"msg": "ok",
		"code": 0,
	}

	if len(data) > 1 { // Almost the length is more than 1, so just check it first.
		respH["data"] = data
	} else if data["data"] != nil {
		respH["data"] = data["data"]
	} else if data != nil && len(data) > 0 {
		respH["data"] = data
	}

	c.Ctx.JSON(http.StatusOK, respH)
}

// error handle
func (c *Context) Error(err error, errCode int) {

	if errCode == 0 {
		errCode = 1
	}

	c.Ctx.JSON(http.StatusOK, gin.H{
		"msg": "error",
		"code": errCode,
		"data": err.Error(),
	})
}

// forbidden handle
func (c *Context) Forbidden() {
	c.Ctx.JSON(http.StatusForbidden, gin.H{
		"msg": "forbidden",
		"code": http.StatusForbidden,
		"data": nil,
	})
	c.Ctx.Abort()
}

// get post data
func (c *Context) getRaw(key string) string {
	res := gjson.GetBytes(c.RawData, key)
	return res.Str
}

// get params
func (c *Context) getParam(key string) string {
	res := c.Ctx.Param(key)
	return res
}