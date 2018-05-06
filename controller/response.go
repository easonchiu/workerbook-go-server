package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

// success
func (r *Response) Success(data gin.H) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"code": 0,
		"data": data,
	})
}

// error
func (r *Response) Error(err error) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"msg": "error",
		"code": 1,
		"data": err.Error(),
	})
}

// forbidden
func (r *Response) Forbidden() {
	r.Ctx.JSON(http.StatusForbidden, gin.H{
		"msg": "forbidden",
		"code": http.StatusForbidden,
		"data": nil,
	})
	r.Ctx.Abort()
}