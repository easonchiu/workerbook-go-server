import { handleActions } from 'easy-action'

const initialState = {
  list: [],
  profile: {},
}

export default handleActions({
  USER_LIST(state, action) {
    return {
      ...state,
      list: action.payload.list,
    }
  },
  USER_PROFILE(state, action) {
    return {
      ...state,
      profile: action.payload,
    }
  }
}, initialState)
