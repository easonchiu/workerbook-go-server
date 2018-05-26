import './style'
import React from 'react'

import AsidePanel from '../asidePanel'

const Item = props => {
  if (props.active) {
    return (
      <p>{props.children}</p>
    )
  }
  return (
    <a href="javascript:;" onClick={props.onClick}>
      {props.children}
    </a>
  )
}

const AsideGroupList = props => {
  const { list, active, itemClick = () => {} } = props
  return (
    <AsidePanel title="成员分组" className="aside-group-list">
      <ul>
        <li>
          <Item active={active === ''} onClick={itemClick.bind(null, '')}>
            <i /><span>全部</span>
          </Item>
        </li>
        {
          list.map(item => {
            return (
              <li key={item.id}>
                <Item active={active === item.id} onClick={itemClick.bind(null, item.id)}>
                  <i />
                  <span>{item.name}</span>
                  <em>{item.count}人</em>
                </Item>
              </li>
            )
          })
        }
      </ul>
    </AsidePanel>
  )
}
export default AsideGroupList