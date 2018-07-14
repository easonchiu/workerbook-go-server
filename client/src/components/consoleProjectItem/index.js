import './style'
import React from 'react'
import ignore from 'src/utils/ignore'

import IconSetting from 'src/components/svg/setting'
import IconAdd from 'src/components/svg/add'
import IconDelete from 'src/components/svg/delete'

const ProjectItem = props => {
  const source = props.source || {}

  const renderTools = () => (
    <div className="tools">
      <a
        href="javascript:;"
        className="del"
        onClick={() => {
          props.onProjectDeleteClick && props.onProjectDeleteClick(props.source)
        }}
      >
        <IconDelete />
      </a>
      <a
        href="javascript:;"
        onClick={() => {
          props.onProjectEditClick && props.onProjectEditClick(props.source)
        }}
      >
        <IconSetting />
      </a>
    </div>
  )

  const renderMissions = () => (
    <div className="missions">
      <h6>包含任务 {source.missions ? source.missions.length + '个' : ''}</h6>
      {
        source.missions && source.missions.length ?
          <ul>
            {
              source.missions.map(item => (
                <li key={item.id}>
                  <p>{item.name}</p>
                  <a
                    href="javascript:;"
                    onClick={() => {
                      props.onMissionEditClick &&
                      props.onMissionEditClick(item, ignore(source, 'missions'))
                    }}
                  >
                    编辑
                  </a>
                </li>
              ))
            }
          </ul> :
          <p className="empty">暂无任务</p>
      }
      <a
        href="javascript:;"
        className="append"
        onClick={() => {
          props.onAppendMissionClick && props.onAppendMissionClick(props.source)
        }}
      >
        <IconAdd />
      </a>
    </div>
  )

  const renderFooter = () => (
    <footer>
      <span>时间周期</span>
      {new Date(source.deadline).format('yyyy年MM月dd日 hh:mm')}
      {' ~ '}
      {new Date(source.deadline).format('MM月dd日 hh:mm')}
    </footer>
  )

  return (
    <div className="wbc-project-item">
      <h2>
        {
          source.weight === 2 ?
            <span className="weight-2">重要</span> :
            source.weight === 3 ?
              <span className="weight-3">紧急</span> :
              null
        }
        <a href="#">{source.name}</a>
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
      {renderMissions()}
      {renderTools()}
      {renderFooter()}
    </div>
  )
}
export default ProjectItem