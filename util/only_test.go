package util

import (
  "fmt"
  "gopkg.in/mgo.v2/bson"
  "testing"
  "time"
)

func TestIsType(t *testing.T) {

  var str interface{} = "string"

  isType(str, TypeString)

  var i interface{} = 123

  isType(i, TypeInt)

  var b interface{} = false

  isType(b, TypeBool)

  var oid interface{} = bson.NewObjectId()

  isType(oid, TypeObjectId)

  var tm interface{} = time.Now()

  isType(tm, TypeTime)

  var bm interface{} = bson.M{}

  isType(bm, TypeBsonM)

  var arr interface{} = []bson.M{
    {"a": 1},
    {"b": 2},
  }

  isType(arr, TypeAny)

}

func TestOnly(t *testing.T) {

  m := bson.M{
    "str":  "string",
    "bool": false,
    "time": time.Now(),
    "id":   bson.NewObjectId(),
    "int":  123,
    "m": bson.M{
      "sub": "test",
    },
    "arr": []bson.M{
      {"a": 1},
      {"b": 2},
    },
  }

  Only(m, Keys{
    "str":  TypeString,
    "bool": TypeBool,
    "time": TypeTime,
    "id":   TypeObjectId,
    "int":  TypeInt,
    "m":    TypeBsonM,
    "arr":  TypeAny,
  })

  fmt.Println(m)

}
