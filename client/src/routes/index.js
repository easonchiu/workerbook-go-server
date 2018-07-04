import APP_CONFIG from '../../appConfig'

import React from 'react'
import AsyncComponent from 'src/hoc/asyncComponent'
import { BrowserRouter, Route, Redirect, Switch } from 'react-router-dom'

const Login = AsyncComponent(() => import('src/views/login'))
const Index = AsyncComponent(() => import('src/views/index'))
const Project = AsyncComponent(() => import('src/views/project'))
const Chart = AsyncComponent(() => import('src/views/chart'))
const ChartProject = AsyncComponent(() => import('src/views/chartProject'))
const ChartGroup = AsyncComponent(() => import('src/views/chartGroup'))


const CreateUser = AsyncComponent(() => import('src/views/createUser'))
const CreateGroup = AsyncComponent(() => import('src/views/createGroup'))
const CreateProject = AsyncComponent(() => import('src/views/createProject'))


// 配置路由
const Routes = () => {
  return (
    <BrowserRouter basename={APP_CONFIG.basename}>
      <Switch>
        <Route exact path="/login" component={Login} />
        <Route exact path="/index" component={Index} />
        <Route exact path="/project" component={Project} />
        <Route exact path="/chart" component={Chart} />
        <Route exact path="/chart/project" component={ChartProject} />
        <Route exact path="/chart/group" component={ChartGroup} />

        <Route exact path="/createUser" component={CreateUser} />
        <Route exact path="/createGroup" component={CreateGroup} />
        <Route exact path="/createProject" component={CreateProject} />
        <Redirect to="/index" />
      </Switch>
    </BrowserRouter>
  )
}

export default Routes
