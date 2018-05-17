import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch user's daily list.
const fetchUsersDailyList = payload => async dispatch => {
  const res =  await http.request({
    url: '/daily',
    method: 'GET',
    data: payload
  })
  dispatch(createAction('DAILY_USERS_DAILY_LIST')(res))
}


export default {
  fetchUsersDailyList,
}
