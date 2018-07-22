package util

import (
  "github.com/gin-gonic/gin"
  "strings"
)

func Forget(m gin.H, key string) {
  key = strings.TrimSpace(key)
  keys := strings.Split(key, " ")
  for _, v := range keys {
    if v != "" {
      delete(m, v)
    }
  }
}

func ForgetArr(arr interface{}, key string) {
  switch arr.(type) {
  case []gin.H:
    for _, item := range arr.([]gin.H) {
      Forget(item, key)
    }
  default:
  }
}