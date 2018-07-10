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
    method: 'PUT',
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
  dispatch(createAction('DEPARTMENT_LIST')(res))
}

// fetch departments list for select.
const fetchSelectList = () => async dispatch => {
  const res = await http.request({
    url: '/departments/all',
    method: 'GET',
  })
  dispatch(createAction('DEPARTMENT_SELECT_LIST')(res))
}


export default {
  create,
  update,
  fetchList,
  fetchSelectList,
}
