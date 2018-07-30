package model

import "github.com/gin-gonic/gin"

const REMEMBER = "_Remember_"

func exist(k string, keys ... string) bool {
  for _, i := range keys {
    if i == k {
      return true
    }
  }
  return false
}

func forget(m gin.H, keys ... string) {
  for _, v := range keys {
    if _, ok := m[v]; ok {
      delete(m, v)
    }
  }
}

func remember(m gin.H, keys ... string) {
  for k := range m {
    if !exist(k, keys...) {
      delete(m, k)
    }
  }
}
