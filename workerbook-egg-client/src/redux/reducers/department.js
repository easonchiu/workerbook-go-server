import { handleActions } from 'easy-action'

const initialState = {
  departments: {
    list: [],
    count: 0,
    skip: 0,
    limit: 0,
  },
  select: {
    list: []
  }
}

export default handleActions({
  DEPARTMENT_LIST(state, action) {
    return {
      ...state,
      departments: {
        list: action.payload.list || [],
        skip: action.payload.skip,
        limit: action.payload.limit,
        count: action.payload.count,
      }
    }
  },
  DEPARTMENT_SELECT_LIST(state, action) {
    return {
      ...state,
      select: {
        list: action.payload || [],
      }
    }
  }
}, initialState)
