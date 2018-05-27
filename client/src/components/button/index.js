import './style'
import React from 'react'
import classNames from 'classnames'

const Button = props => {
  const css = classNames('x-button', props.className)

  return (
    <button {...props} className={css}>
      {props.children}
    </button>
  )
}
export default Button