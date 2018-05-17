import APP_CONFIG from '../../appConfig'

import React from 'react'
import AsyncComponent from 'src/hoc/asyncComponent'
import { BrowserRouter, Route, Redirect, Switch } from 'react-router-dom'

const Login = AsyncComponent(() => import('src/views/login'))
const Index = AsyncComponent(() => import('src/views/index'))

// 配置路由
const Routes = () => {
  return (
    <BrowserRouter basename={APP_CONFIG.basename}>
      <Switch>
        <Route exact path="/login" component={Login} />
        <Route exact path="/index" component={Index} />
        <Redirect to="/index" />
      </Switch>
    </BrowserRouter>
  )
}

export default Routes
