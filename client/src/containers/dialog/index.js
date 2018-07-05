import './style'
import React from 'react'
import { createPortal } from 'react-dom'
import classNames from 'classnames'

class Dialog extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      visible: props.visible,
      ani: 'out'
    }

    if (props.visible) {
      setTimeout(() => {
        this.setState({
          ani: 'in'
        })
      })
    }

    this.el = document.getElementById('app_dialogs')
    if (!this.el) {
      this.el = document.createElement('div')
      this.el.id = 'app_dialogs'
      document.body.appendChild(this.el)
    }
    this.dialogEl = document.createElement('div')
    this.el.appendChild(this.dialogEl)
  }

  componentWillUnmount() {
    this.el.removeChild(this.dialogEl)
  }

  static getDerivedStateFromProps(np, ps) {
    if (np.visible) {
      return {
        visible: true
      }
    }
    else {
      return {
        ani: 'out'
      }
    }
  }

  componentDidUpdate(pp, ps) {
    if (!ps.visible && this.state.visible) {
      setTimeout(() => {
        this.setState({
          ani: 'in'
        })
        this.props.onStatusChange && this.props.onStatusChange(true)
      })
    }
    else if (ps.ani === 'in' && this.state.ani === 'out') {
      this.props.onStatusChange && this.props.onStatusChange(false)
      setTimeout(() => {
        this.setState({
          visible: false
        })
      }, 250)
    }
  }

  renderContent() {
    const css = classNames('app-dialog', this.props.className, `app-dialog--${this.state.ani}`)
    return (
      <div
        style={{ display: this.state.visible ? '' : 'none' }}
        className={css}
      >
        <div className="app-dialog__content">
          {this.props.children}
        </div>
        <div className="app-dialog__bg" />
      </div>
    )
  }

  render() {
    return createPortal(
      this.renderContent(),
      this.dialogEl
    )
  }
}

export default Dialog