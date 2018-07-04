import './style'
import React from 'react'

import classNames from 'classnames'

const UserHeader = props => {
  const { name } = props
  const css = classNames('user-header', {
    mini: props.mini
  })
  return (
    <div className={css}>
      {name}
    </div>
  )
}
export default UserHeader