import http from 'src/utils/http'
import { createAction } from 'easy-action'
import ignore from 'src/utils/ignore'

// user login
const login = payload => async () => {
  const res = await http.request({
    url: '/login',
    method: 'POST',
    data: payload,
  })
  return res
}

// my daily
const fetchMyTodayDaily = () => async dispatch => {
  const res = await http.request({
    url: '/users/my/dailies/today',
    method: 'GET',
  })
  dispatch(createAction('USER_MY_TODAY_DAILY')(res))
}


// append daily item.
const appendDailyItem = ({ record, progress, project }) => async () => {
  const res = await http.request({
    url: '/users/my/dailies/today/items',
    method: 'POST',
    data: {
      record: record.trim(),
      progress,
      project: project.trim(),
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
const fetchProfile = () => async (dispatch, getState) => {
  const state = getState()
  if (state.user$.profile.id) {
    return
  }
  const res = await http.request({
    url: '/profile',
    method: 'GET',
  })
  dispatch(createAction('USER_PROFILE')(res))
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

// update user
const update = payload => async () => {
  const res = await http.request({
    url: '/users/' + payload.id,
    method: 'PUT',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch users list.
const fetchList = ({ departmentId, skip, limit = 10 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/users',
    method: 'GET',
    params: {
      departmentId,
      skip,
      limit,
    }
  })
  dispatch(createAction('USER_LIST')(res))
}

export default {
  login,
  create,
  update,
  fetchList,
  fetchProfile,
  appendDailyItem,
  deleteDailyItem,
  fetchMyTodayDaily,
}
