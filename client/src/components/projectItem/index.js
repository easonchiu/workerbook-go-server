import './style'
import React from 'react'

import Button from 'src/components/button'

const ProjectItem = props => {
  const source = props.source || {}
  return (
    <div className="project-item">
      <h2>
        {
          source.weight === 2 ?
            <span className="weight-2">重要</span> :
            source.weight === 3 ?
              <span className="weight-3">紧急</span> :
              null
        }
        {source.name}
      </h2>
      <div className="departments">
        <span>参与部门</span>
        {
          source.departments ?
            source.departments.map(item => (
              <em key={item.id}>{item.name}</em>
            )) :
            null
        }
      </div>
      <footer>
        <span>时间周期</span>
        {new Date(source.createTime).format('yyyy年MM月dd日 hh:mm')}
        {' ~ '}
        {new Date(source.deadline).format('MM月dd日 hh:mm')}
      </footer>
    </div>
  )
}
export default ProjectItem