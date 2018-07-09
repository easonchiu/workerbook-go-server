import { createAction } from 'easy-action'
import http from 'src/utils/http'
import ignore from 'src/utils/ignore'

// create user
const createUser = payload => async () => {
  const res = await http.request({
    url: '/console/users',
    method: 'POST',
    data: payload,
  })
  return res
}

// update user
const updateUser = payload => async () => {
  const res = await http.request({
    url: '/console/users/' + payload.id,
    method: 'PUT',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch users list.
const fetchUsersList = ({ departmentId, skip, limit = 10 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/console/users',
    method: 'GET',
    params: {
      departmentId,
      skip,
      limit,
    }
  })
  dispatch(createAction('CONSOLE_USERS_LIST')(res))
}

// create user
const createDepartment = payload => async () => {
  const res = await http.request({
    url: '/console/departments',
    method: 'POST',
    data: payload,
  })
  return res
}

// update user
const updateDepartment = payload => async () => {
  const res = await http.request({
    url: '/console/departments/' + payload.id,
    method: 'PUT',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch users list.
const fetchDepartmentsList = ({ skip, limit = 10 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/console/departments',
    method: 'GET',
    params: {
      skip,
      limit,
    }
  })
  dispatch(createAction('CONSOLE_DEPARTMENTS_LIST')(res))
}

export default {
  createUser,
  updateUser,
  fetchUsersList,

  createDepartment,
  updateDepartment,
  fetchDepartmentsList,
}
