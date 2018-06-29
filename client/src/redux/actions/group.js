import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch group list.
const fetchList = ({ skip = 0, limit = 4 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/groups',
    method: 'GET',
    params: {
      skip,
      limit,
    }
  })
  dispatch(createAction('GROUP_LIST')(res))
}

// create group
const create = payload => async () => {
  const res = await http.request({
    url: '/groups',
    method: 'POST',
    data: payload
  })
  return res
}

export default {
  fetchList,
  create,
}
