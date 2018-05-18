
// reset css
import 'src/assets/css/reset'

// some utils
import 'src/utils/dateFormat'

// base framework
import React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'

// store
import configureStore from 'src/redux/store'
const store = configureStore()

// routes
import Routers from 'src/routes'

// render to #root
render(
  <Provider store={store}>
    <Routers />
  </Provider>,
  document.getElementById('root'),
)