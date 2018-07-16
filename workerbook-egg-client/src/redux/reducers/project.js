import { handleActions } from 'easy-action'

const initialState = {
  projects: {
    list: [],
    count: 0,
    skip: 0,
    limit: 0,
  },
}

export default handleActions({
  PROJECT_LIST(state, action) {
    return {
      ...state,
      projects: {
        list: action.payload.list || [],
        count: action.payload.count,
        skip: action.payload.skip,
        limit: action.payload.limit,
      },
    }
  }
}, initialState)
