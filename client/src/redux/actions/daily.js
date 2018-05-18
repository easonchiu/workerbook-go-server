import { createAction } from 'easy-action'
import http from 'src/utils/http'

// fetch daily list by day.
const fetchDailyListByDay = payload => async dispatch => {
  const res =  await http.request({
    url: '/daily',
    method: 'GET',
    data: payload
  })
  dispatch(createAction('DAILY_LIST_BY_DAY')(res))
}


export default {
  fetchDailyListByDay,
}
