import { handleActions } from 'easy-action'

const initialState = {
  users: {
    list: [],
    skip: 0,
    limit: 0,
    count: 0,
  },
  departments: {
    list: [],
    skip: 0,
    limit: 0,
    count: 0,
  }
}

export default handleActions({
  CONSOLE_USERS_LIST(state, action) {
    return {
      ...state,
      users: {
        list: action.payload.list || [],
        skip: action.payload.skip,
        limit: action.payload.limit,
        count: action.payload.count,
      }
    }
  },
  CONSOLE_DEPARTMENTS_LIST(state, action) {
    return {
      ...state,
      departments: {
        list: action.payload.list || [],
        skip: action.payload.skip,
        limit: action.payload.limit,
        count: action.payload.count,
      }
    }
  }
}, initialState)
