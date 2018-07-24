package util

import (
  "github.com/gin-gonic/gin"
)

func Forget(m gin.H, keys ... string) {
  for _, v := range keys {
    if _, ok := m[v]; ok {
      delete(m, v)
    }
  }
}

func ForgetArr(arr []gin.H, keys ... string) {
  for _, item := range arr {
    Forget(item, keys...)
  }
}
