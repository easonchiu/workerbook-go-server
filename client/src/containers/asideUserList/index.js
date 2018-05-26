import './style'
import React from 'react'

import AsidePanel from '../asidePanel'
import UserHeader from 'src/components/userHeader'

const AsideUserList = props => {
  const { list } = props
  return (
    <AsidePanel title="全部成员" className="aside-user-list">
      <div className="user-list">
        {
          list.map(item => {
            return (
              <UserHeader name={item.nickname} key={item.id} />
            )
          })
        }
      </div>
    </AsidePanel>
  )
}
export default AsideUserList