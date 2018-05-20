import { handleActions } from 'easy-action'

const initialState = {
  mine: [],
  list: [],
}

export default handleActions({
  DAILY_MY(state, action) {
    return {
      ...state,
      mine: action.payload,
    }
  },
  DAILY_LIST_BY_DAY(state, action) {
    return {
      ...state,
      list: action.payload.list,
    }
  }
}, initialState)
