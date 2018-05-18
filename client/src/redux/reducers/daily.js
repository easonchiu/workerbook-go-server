import { handleActions } from 'easy-action'

const initialState = {
  list: [],
}

export default handleActions({
  DAILY_LIST_BY_DAY(state, action) {
    return {
      ...state,
      list: action.payload.list,
    }
  }
}, initialState)
