package service

import (
	"github.com/gin-gonic/gin"
	"errors"
)

func GetUserInfoById(id string) (gin.H, error) {
	if id != "123" {
		return nil, errors.New("cant find")
	}

	return gin.H{
		"id": id,
		"name": "eason.chiu",
		"age": 18,
	}, nil
}