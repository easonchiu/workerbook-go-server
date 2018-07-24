package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

type test struct {
  name string
}

type testList struct {
  list  []test
  count int
}

func (t test) getMap() gin.H {
  return gin.H{
    "name": t.name,
  }
}

func getList() *testList {
  res := new(testList)

  for i := 0; i < 10; i++ {
    res.list = append(res.list, test{
      name: "hello",
    })
  }

  res.count = 10

  return res
}

func (t testList) each(i func(test) gin.H) gin.H {
  data := gin.H{
    "count": t.count,
  }

  var list []gin.H
  for _, item := range t.list {
    list = append(list, i(item))
  }

  data["list"] = list

  return data
}

func main() {
  l := getList()

  data := l.each(func(i test) gin.H {
    d := i.getMap()
    d["append"] = true
    return d
  })

  fmt.Println(data)
}
