import { combineReducers } from 'redux'

import user$ from './user'
import daily$ from './daily'
import group$ from './group'

export default combineReducers({
  user$,
  daily$,
  group$,
})
