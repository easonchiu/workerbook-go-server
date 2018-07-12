import { handleActions } from 'easy-action'

const initialState = {
  profile: {},
  users: {
    list: [],
    departmentId: null,
    count: 0,
    skip: 0,
    limit: 0,
  },
  one: {}
}

export default handleActions({
  USER_MY_TODAY_DAILY(state, action) {
    return {
      ...state,
      fetchMyTodayDaily: action.payload,
    }
  },
  USER_LIST(state, action) {
    return {
      ...state,
      users: {
        list: action.payload.list,
        departmentId: action.payload.departmentId,
        count: action.payload.count,
        skip: action.payload.skip,
        limit: action.payload.limit,
      },
    }
  },
  USER_PROFILE(state, action) {
    return {
      ...state,
      profile: action.payload,
    }
  }
}, initialState)
