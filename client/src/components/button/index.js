import './style'
import React from 'react'
import classNames from 'classnames'

const Button = props => {
  const css = classNames('x-button', {
    mini: props.mini,
    light: props.light,
  }, props.className)

  return (
    <button className={css}>
      {props.children}
    </button>
  )
}
export default Button