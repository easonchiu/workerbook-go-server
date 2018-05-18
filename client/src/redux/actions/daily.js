import { createAction } from 'easy-action'
import http from 'src/utils/http'

// append daily.
const append = ({ record, project }) => async () => {
  const res = await http.request({
    url: '/daily',
    method: 'POST',
    data: {
      record,
      progress: 50,
      pid: project,
    }
  })
  return res
}

// fetch daily list by day.
const fetchDailyListByDay = ({ skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/daily',
    method: 'GET',
    params: {
      skip,
      limit,
    }
  })
  dispatch(createAction('DAILY_LIST_BY_DAY')(res))
}


export default {
  append,
  fetchDailyListByDay,
}
