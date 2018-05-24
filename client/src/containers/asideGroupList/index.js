import './style'
import React from 'react'

import AsidePanel from '../asidePanel'

const AsideGroupList = props => {
  const { list } = props
  return (
    <AsidePanel title="成员分组" className="aside-group-list">
      <ul>
        <li>
          <a href="javascript:;" onClick={props.itemClick.bind(null, '')}>
            全部
          </a>
        </li>
        {
          list.map(item => {
            return (
              <li key={item.id}>
                <a href="javascript:;" onClick={props.itemClick.bind(null, item.id)}>
                  {item.name}
                  <span>{item.count}</span>
                </a>
              </li>
            )
          })
        }
      </ul>
    </AsidePanel>
  )
}
export default AsideGroupList