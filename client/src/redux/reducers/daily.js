import { handleActions } from 'easy-action'

const initialState = {
  list: [],
}

export default handleActions({
  DAILY_USERS_DAILY_LIST(state, action) {
    return Object.assign({}, state, {
      list: action.payload.list,
    })
  }
}, initialState)
