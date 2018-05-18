import APP_CONFIG from '../../appConfig'

import React from 'react'
import AsyncComponent from 'src/hoc/asyncComponent'
import { BrowserRouter, Route, Redirect, Switch } from 'react-router-dom'

const Login = AsyncComponent(() => import('src/views/login'))
const Index = AsyncComponent(() => import('src/views/index'))
const CreateUser = AsyncComponent(() => import('src/views/createUser'))
const CreateGroup = AsyncComponent(() => import('src/views/createGroup'))

// 配置路由
const Routes = () => {
  return (
    <BrowserRouter basename={APP_CONFIG.basename}>
      <Switch>
        <Route exact path="/login" component={Login} />
        <Route exact path="/index" component={Index} />
        <Route exact path="/createUser" component={CreateUser} />
        <Route exact path="/createGroup" component={CreateGroup} />
        <Redirect to="/index" />
      </Switch>
    </BrowserRouter>
  )
}

export default Routes
