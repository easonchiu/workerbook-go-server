import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch group list.
const fetchList = payload => async dispatch => {
  const res =  await http.request({
    url: '/group',
    method: 'GET',
    data: payload
  })
  dispatch(createAction('GROUP_LIST')(res))
}

// create group
const create = payload => async () => {
  const res =  await http.request({
    url: '/group',
    method: 'POST',
    data: payload
  })
  return res
}

export default {
  fetchList,
  create,
}
