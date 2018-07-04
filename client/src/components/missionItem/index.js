import './style'
import React from 'react'
import classNames from 'classnames'

import UserHeader from 'src/components/userHeader'
import Button from 'src/components/button'

const MissionItem = props => {
  const css = classNames('mission-item', {
    'show-joined': props.showJoined
  })
  return (
    <div className={css}>
      <h2>前端页面开发</h2>
      <p className="project">所属项目：世界杯活动页面开发</p>
      <time>项目截至时间：2018年3月3日</time>
      <div className="tools">
        <Button mini light>项目说明</Button>
        <Button mini light>任务说明</Button>
      </div>
      <div className="progress"><span>45</span></div>
      {
        props.showJoined ?
          <div className="joined-list">
            <UserHeader name="Eason.Chiu" mini to={1} />
            <UserHeader name="牛哥牛哥" mini to={1} />
            <UserHeader name="张小三" mini to={1} />
            <UserHeader name="李四" mini to={1} />
            <UserHeader name="龙五" mini to={1} />
            <p className="more">等32人</p>
            {
              !props.joined ?
                <Button mini className="join">分配</Button> :
                null
            }
          </div> :
          null
      }
    </div>
  )
}
export default MissionItem