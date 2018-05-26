import './style'
import React from 'react'

const UserHeader = props => {
  const { name } = props
  return (
    <div className="user-header">
      {name}
    </div>
  )
}
export default UserHeader