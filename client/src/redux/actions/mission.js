import { createAction } from 'easy-action'
import http from 'src/utils/http'
import ignore from 'src/utils/ignore'

// create mission
const create = payload => async () => {
  const res = await http.request({
    url: '/missions',
    method: 'POST',
    data: payload
  })
  return res
}

// update mission
const update = payload => async () => {
  const res = await http.request({
    url: '/missions/' + payload.id,
    method: 'PUT',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch mission list.
const fetchList = ({ status, skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/missions',
    method: 'GET',
    params: {
      status,
      skip,
      limit,
    }
  })
  dispatch(createAction('MISSION_LIST')(res))
}

// fetch mission one by id
const fetchOneById = id => async dispatch => {
  const res = await http.request({
    url: '/missions/' + id,
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
