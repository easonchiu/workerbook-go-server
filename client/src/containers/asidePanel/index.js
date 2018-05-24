import './style'
import React from 'react'
import cn from 'classnames'

const AsidePanel = props => {
  const { title } = props
  return (
    <div className={cn('aside-panel', props.className)}>
      <h2>{title}</h2>
      {props.children}
    </div>
  )
}
export default AsidePanel