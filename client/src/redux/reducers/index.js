import { combineReducers } from 'redux'

import user$ from './user'
import daily$ from './daily'

export default combineReducers({
  user$,
  daily$,
})
