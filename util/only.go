package util

import (
  "gopkg.in/mgo.v2/bson"
  "reflect"
  "time"
)

type T string

const (
  TypeString   T = "string"
  TypeInt      T = "int"
  TypeBool     T = "bool"
  TypeTime     T = "time.Time"
  TypeObjectId T = "bson.ObjectId"
  TypeBsonM    T = "bson.M"
  TypeAny      T = "*"
)

type Keys map[string]T

func isType(key interface{}, t T) bool {
  if t == TypeAny {
    return true
  }
  return reflect.TypeOf(key).String() == reflect.ValueOf(t).String()
}

func exist(t Keys, k string, v interface{}) bool {
  for key, types := range t {
    if k == key && isType(v, types) {
      return true
    }
  }
  return false
}

// bson.M只允许部分key，其他都删除
func Only(m bson.M, t Keys) {
  for k, v := range m {
    b := exist(t, k, v)
    if !b {
      delete(m, k)
    }
  }
}

// 从bson.M中获取string类型的值
func GetString(m bson.M, key string) (string, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeString) {
      return res.(string), true
    }
    return "", false
  }
  return "", false
}

// 从bson.M中获取int类型的值
func GetInt(m bson.M, key string) (int, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeInt) {
      return res.(int), true
    }
    return 0, false
  }
  return 0, false
}

// 从bson.M中获取time类型的值
func GetTime(m bson.M, key string) (*time.Time, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeTime) {
      t := res.(time.Time)
      return &t, true
    }
    return nil, false
  }
  return nil, false
}

// 从bson.M中获取bool类型的值
func GetBool(m bson.M, key string) (bool, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeBool) {
      return res.(bool), true
    }
    return false, false
  }
  return false, false
}

// 从bson.M中获取objectId类型的值
func GetObjectId(m bson.M, key string) (*bson.ObjectId, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeObjectId) {
      id := res.(bson.ObjectId)
      return &id, true
    }
    return nil, false
  }
  return nil, false
}

// 从bson.M中获取bson.M类型的值
func GetBsonM(m bson.M, key string) (bson.M, bool) {
  if res, ok := m[key]; ok {
    if isType(res, TypeObjectId) {
      return res.(bson.M), true
    }
    return nil, false
  }
  return nil, false
}

// 从bson.M中获取interface{}类型的值
func GetAny(m bson.M, key string) (interface{}, bool) {
  if res, ok := m[key]; ok {
    return res, true
  }
  return nil, false
}
