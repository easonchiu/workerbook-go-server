import { createAction } from 'easy-action'
import http from 'src/utils/http'

// my daily
const mine = () => async dispatch => {
  const res = await http.request({
    url: '/users/my/dailies/today',
    method: 'GET',
  })
  dispatch(createAction('DAILY_MY')(res))
}

// fetch daily list by day.
const fetchListByDay = ({ skip, limit } = {}) => async dispatch => {
  const res = await http.request({
    url: '/dailies',
    method: 'GET',
    params: {
      skip,
      limit,
    }
  })
  dispatch(createAction('DAILY_LIST_BY_DAY')(res))
}

export default {
  mine,
  fetchListByDay,
}
