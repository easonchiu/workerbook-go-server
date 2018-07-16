import { createAction } from 'easy-action'
import http from 'src/utils/http'
import ignore from 'src/utils/ignore'

// create department
const create = payload => async () => {
  const res = await http.request({
    url: '/departments',
    method: 'POST',
    data: payload,
  })
  return res
}

// update department
const update = payload => async () => {
  const res = await http.request({
    url: '/departments/' + payload.id,
    method: 'PATCH',
    data: ignore(payload, 'id'),
  })
  return res
}

// fetch departments list.
const fetchList = ({ skip = 0, limit = 10 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/departments',
    method: 'GET',
    params: {
      skip,
      limit,
    }
  })
  const data = {
    list: res,
    count: res.length,
    skip,
    limit,
  }
  dispatch(createAction('DEPARTMENT_LIST')(data))
}

// fetch departments list for select.
const fetchSelectList = () => async dispatch => {
  const res = await http.request({
    url: '/departments',
    method: 'GET',
  })
  dispatch(createAction('DEPARTMENT_SELECT_LIST')(res))
}

// fetch department one by id
const fetchOneById = id => async dispatch => {
  const res = await http.request({
    url: '/departments/' + id,
    method: 'GET',
  })
  return res
}

// delete department
const del = id => async dispatch => {
  const res = await http.request({
    url: '/departments/' + id,
    method: 'DELETE',
  })
  return res
}

export default {
  create,
  update,
  del,
  fetchList,
  fetchOneById,
  fetchSelectList,
}
