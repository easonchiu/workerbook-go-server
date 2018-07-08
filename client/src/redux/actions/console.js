import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch users list.
const fetchUsersList = ({ departmentId, skip, limit = 10 } = {}) => async dispatch => {
  const res = await http.request({
    url: '/console/users',
    method: 'GET',
    params: {
      departmentId,
      skip,
      limit,
    }
  })
  res.skip = skip
  res.limit = limit
  dispatch(createAction('CONSOLE_USERS_LIST')(res))
}

export default {
  fetchUsersList,
}
