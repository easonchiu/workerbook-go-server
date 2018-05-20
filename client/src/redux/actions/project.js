import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch group list.
const fetchList = ({ status, skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/projects',
    method: 'GET',
    params: {
      status,
      skip,
      limit,
    }
  })
  dispatch(createAction('PROJECT_LIST')(res))
}

// create group
const create = payload => async () => {
  const res = await http.request({
    url: '/projects',
    method: 'POST',
    data: payload
  })
  return res
}

export default {
  fetchList,
  create,
}
