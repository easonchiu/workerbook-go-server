import { handleActions } from 'easy-action'

const initialState = {
  fetchMyTodayDaily: [],
  list: [],
  activeGroup: '', // 当前list对应的id
  profile: {},
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
      list: action.payload.list,
      activeGroup: action.payload.gid || '',
    }
  },
  USER_PROFILE(state, action) {
    return {
      ...state,
      profile: action.payload,
    }
  }
}, initialState)
