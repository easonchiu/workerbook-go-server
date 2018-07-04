import './style'
import React from 'react'
import { createPortal } from 'react-dom'

class Dialog extends React.PureComponent {
  constructor(props) {
    super(props)
    this.el = document.createElement('div')
  }

  render() {
    return createPortal(
      <div>123</div>,
      this.el
    )
  }
}

export default Dialog