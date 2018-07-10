import './style'
import React, { cloneElement } from 'react'
import classNames from 'classnames'

class Select extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      visible: false
    }
  }

  listener = e => {
    document.removeEventListener('click', this.listener)
    this.setState({
      visible: false
    })
  }

  componentWillUnmount() {
    document.removeEventListener('click', this.listener)
  }

  transVisible = () => {
    const nextVisible = !this.state.visible
    this.setState({
      visible: nextVisible
    })
    if (nextVisible) {
      document.addEventListener('click', this.listener)
    }
    else {
      document.removeEventListener('click', this.listener)
    }
  }

  render() {
    const { props } = this
    const current = {}
    let children = this.props.children

    if (props.children && !Array.isArray(props.children)) {
      children = [props.children]
    }

    children = children ?
      children.map((item, index) => {
        if (props.value === item.props.value) {
          current.value = item.props.value
          current.text = item.props.children
        }
        return cloneElement(item, {
          value: item.props.value || '',
          key: index,
          onClick: props.onClick,
          current: props.value,
        })
      }) :
      <Option
        className="wb-select__list-none"
        value=""
        onClick={props.onClick}
      >
        请选择
      </Option>

    return (
      <div className="wb-select">
        <a
          href="javascript:;"
          className="wb-select__value"
          onClick={this.transVisible}
        >
          {current.text}
        </a>
        {
          this.state.visible ?
            <ul className="wb-select__list">
              {children}
            </ul> :
            null
        }
      </div>
    )
  }
}

const Option = props => {
  return (
    <li>
      <a
        href="javascript:;"
        name={props.value}
        onClick={() => {
          props.onClick && props.onClick(props.value)
        }}
        className={classNames({
          'wb-select__list-active': props.value === props.current
        }, props.className)}
      >
        {props.children}
      </a>
    </li>
  )
}

Select.Option = Option

export default Select