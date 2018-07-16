import './style'
import React from 'react'

import Button from 'src/components/button'

const ProjectItem = props => {
  const source = props.source || {}
  return (
    <div className="project-item">
      <h2>{source.name}</h2>
      <time>截至日期：{new Date(source.deadline).format('yyyy年MM月dd日 hh时mm分')}</time>
      <p className="join">
        参与部门：
        {
          source.departments && (source.departments.map(item => item.name)).join('、')
        }
      </p>
      <a className="mission-count" onClick={props.onMissionClick}>
        <p>任务数</p>
        <span>{source.missionCount}</span>
      </a>
      <div className="tools">
        <Button mini light>项目说明</Button>
      </div>
      <div className="progress"><span>{source.progress}</span></div>
    </div>
  )
}
export default ProjectItem