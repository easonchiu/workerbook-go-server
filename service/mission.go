package service

import (
  "errors"
  "time"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/mongo"
)

// 创建任务
func CreateMission(data model.Mission) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrMissionNameEmpty)
  errgo.ErrorIfLenLessThen(data.Name, 4, errgo.ErrMissionNameTooShort)
  errgo.ErrorIfLenMoreThen(data.Name, 15, errgo.ErrMissionNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrMissionDeadlineTooLate)
  errgo.ErrorIfLenMoreThen(data.Description, 500, errgo.ErrMissionDescriptionTooLong)

  if err = errgo.PopError(); err != nil {
    return err
  }

  // insert it.
  err = db.C(model.MissionCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  return nil
}
