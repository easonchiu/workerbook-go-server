import './style'
import React from 'react'

import AsidePanel from '../asidePanel'

const AsideProjectList = props => {
  const { list } = props
  return (
    <AsidePanel title="项目分类" className="aside-project-list">
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
                </a>
              </li>
            )
          })
        }
      </ul>
    </AsidePanel>
  )
}
export default AsideProjectList