import { combineReducers } from 'redux'

import user$ from './user'
import daily$ from './daily'
import group$ from './group'
import project$ from './project'

export default combineReducers({
  user$,
  daily$,
  group$,
  project$,
})
