import './style'
import React from 'react'
import { Link } from 'react-router-dom'

const Header = props => {
  const { profile } = props
  return (
    <header className="app-header">
      <div className="app-header__inner">
        {profile.nickname} - {profile.groupName}
        <Link to="/createGroup">Create Group</Link>
        {' | '}
        <Link to="/createUser">Create User</Link>
        {' | '}
        <Link to="/createProject">Create Project</Link>
      </div>
    </header>
  )
}

const Body = props => {
  return (
    <div className="app-body">
      {props.children}
    </div>
  )
}

Body.Main = props => {
  return (
    <main className="app-body__main">
      {props.children}
    </main>
  )
}

Body.Aside = props => {
  return (
    <aside className="app-body__aside">
      {props.children}
    </aside>
  )
}

export default {
  Header,
  Body,
}