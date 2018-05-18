// import { createAction } from 'easy-action'
import http from 'src/utils/http'

// user login
const login = payload => async () => {
  const res =  await http.request({
    url: '/user/login',
    method: 'POST',
    data: payload
  })
  return res
}

// create user
const create = payload => async () => {
  const res =  await http.request({
    url: '/user',
    method: 'POST',
    data: payload
  })
  return res
}


export default {
  login,
  create,
}
