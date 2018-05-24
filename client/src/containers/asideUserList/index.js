import './style'
import React from 'react'

import AsidePanel from '../asidePanel'

const AsideUserList = props => {
  const { list } = props
  return (
    <AsidePanel title="全部成员" className="aside-user-list">
      <ul>
        {
          list.map(item => {
            return (
              <li key={item.id}>
                <a href="javascript:;" onClick={props.itemClick.bind(null, item.id)}>
                  {item.nickname}
                </a>
              </li>
            )
          })
        }
      </ul>
    </AsidePanel>
  )
}
export default AsideUserList