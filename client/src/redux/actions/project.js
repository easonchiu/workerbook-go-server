import { createAction } from 'easy-action'
import http from 'src/utils/http'
import ignore from 'src/utils/ignore'

// create project
const create = payload => async () => {
  const res = await http.request({
    url: '/projects',
    method: 'POST',
    data: payload
  })
  return res
}

// update project
const update = payload => async () => {
  const res = await http.request({
    url: '/projects/' + payload.id,
    method: 'PUT',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch project list.
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

// fetch project one by id
const fetchOneById = id => async dispatch => {
  const res = await http.request({
    url: '/projects/' + id,
    method: 'GET',
  })
  return res
}

export default {
  create,
  update,
  fetchOneById,
  fetchList,
}
