import { handleActions } from 'easy-action'

const initialState = {
  list: [],
}

export default handleActions({
  GROUP_LIST(state, action) {
    return {
      ...state,
      list: action.payload.list,
    }
  }
}, initialState)
