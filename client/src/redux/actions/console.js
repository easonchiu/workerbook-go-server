import { createAction } from 'easy-action'
import http from 'src/utils/http'

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
    url: '/console/users',
    method: 'PUT',
    data: payload,
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

export default {
  createUser,
  updateUser,
  fetchUsersList,
}
