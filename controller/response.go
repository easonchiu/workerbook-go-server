package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ctx *gin.Context
}

func (r *Response) Success(data gin.H) {
	r.ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"code": 0,
		"data": data,
	})
}

func (r *Response) Error(err error) {
	r.ctx.JSON(http.StatusOK, gin.H{
		"msg": "error",
		"code": 1,
		"data": err.Error(),
	})
}
