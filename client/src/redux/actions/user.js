import http from 'src/utils/http'
import { createAction } from 'easy-action'

// user login
const login = payload => async () => {
  const res = await http.request({
    url: '/login',
    method: 'POST',
    data: payload,
  })
  return res
}

// create user
const create = payload => async () => {
  const res = await http.request({
    url: '/users',
    method: 'POST',
    data: payload,
  })
  return res
}

// append daily item.
const appendDailyItem = ({ record, progress, project }) => async () => {
  const res = await http.request({
    url: '/users/my/dailies/today/items',
    method: 'POST',
    data: {
      record,
      progress,
      project,
    }
  })
  return res
}

// delete daily item
const deleteDailyItem = ({ id }) => async () => {
  const res = await http.request({
    url: '/users/my/dailies/today/items/' + id,
    method: 'DELETE',
  })
  return res
}

// my profile
const myProfile = () => async dispatch => {
  const res = await http.request({
    url: '/users/my',
    method: 'GET',
  })
  dispatch(createAction('USER_PROFILE')(res))
}

// fetch list
const fetchList = ({ gid, skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/users',
    method: 'GET',
    params: {
      gid,
      skip,
      limit,
    }
  })
  dispatch(createAction('USER_LIST')(res))
}

export default {
  login,
  create,
  fetchList,
  myProfile,
  appendDailyItem,
  deleteDailyItem,
}
