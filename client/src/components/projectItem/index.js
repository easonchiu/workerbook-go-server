import './style'
import React from 'react'

import Button from 'src/components/button'

const ProjectItem = props => {
  return (
    <div className="project-item">
      <h2>世界杯活动页面开发</h2>
      <time>截至日期：2018年3月3日</time>
      <p className="join">参与部门：系统开发组、测试组</p>
      <a className="mission-count" onClick={props.onMissionClick}>
        <p>任务数</p>
        <span>4</span>
      </a>
      <div className="tools">
        <Button mini>项目说明</Button>
      </div>
    </div>
  )
}
export default ProjectItem