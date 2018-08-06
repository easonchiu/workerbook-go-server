package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "sort"
  "strconv"
)

func main() {

  data := []gin.H{
    {"day": "20180605"},
    {"day": "20180602"},
    {"day": "20180607"},
    {"day": "20180601"},
    {"day": "20180604"},
    {"day": "20180606"},
    {"day": "20180603"},
    {"day": "20180609"},
    {"day": "201806012"},
    {"day": "201806010"},
    {"day": "20180608"},
    {"day": "201806011"},
  }

  sort.Slice(data, func(i, j int) bool {
    d1 := data[i]["day"]
    d2 := data[j]["day"]

    sd1, _ := strconv.Atoi(d1.(string))
    sd2, _ := strconv.Atoi(d2.(string))

    return sd1 < sd2
  })

  fmt.Println(data)

}
