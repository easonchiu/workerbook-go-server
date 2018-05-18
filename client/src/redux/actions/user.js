import http from 'src/utils/http'
import { getToken } from 'src/utils/token'
import { createAction } from 'easy-action'

// user login
const login = payload => async () => {
  const res = await http.request({
    url: '/user/login',
    method: 'POST',
    data: payload,
  })
  return res
}

// create user
const create = payload => async () => {
  const res = await http.request({
    url: '/user',
    method: 'POST',
    data: payload,
  })
  return res
}

// my profile
const myProfile = () => async dispatch => {
  const token = getToken() || '0'
  const res = await http.request({
    url: '/user/' + token,
    method: 'GET',
  })
  dispatch(createAction('USER_PROFILE')(res))
}

// fetch list
const fetchList = ({ gid, skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/user',
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
}
