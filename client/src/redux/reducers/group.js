import { handleActions } from 'easy-action'

const initialState = {
  list: {
    data: [],
    count: 0,
    skip: 0,
    limit: 0,
  },
}

export default handleActions({
  GROUP_LIST(state, action) {
    return {
      ...state,
      list: {
        data: action.payload.list,
        skip: action.payload.skip,
        limit: action.payload.limit,
        count: action.payload.count,
      }
    }
  }
}, initialState)
