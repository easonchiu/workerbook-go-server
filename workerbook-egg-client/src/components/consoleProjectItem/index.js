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
      <IconDelete.A
        className="del"
        onClick={() => {
          props.onDelProjectClick && props.onDelProjectClick(props.source)
        }}
      />
      <IconSetting.A
        onClick={() => {
          props.onEditProjectClick && props.onEditProjectClick(props.source)
        }}
      />
    </div>
  )

  const renderMissions = () => {
    const Item = data => {
      return (
        <li key={data.id}>
          <header>
            <h6>{data.name}</h6>
            <div className="tools">
              <a
                href="javascript:;"
                onClick={() => {
                  props.onDelMissionClick &&
                  props.onDelMissionClick(data, ignore(source, 'missions'))
                }}
              >
                删除
              </a>
              <a
                href="javascript:;"
                onClick={() => {
                  props.onEditMissionClick &&
                  props.onEditMissionClick(data, ignore(source, 'missions'))
                }}
              >
                编辑
              </a>
            </div>
          </header>
          <p>
            <strong>[{(new Date(data.deadline)).format('MM-dd hh:mm')}]</strong>
            {data.description}
          </p>
        </li>
      )
    }
    return (
      <div className="missions">
        <h5>包含任务 {source.missions ? source.missions.length + '个' : ''}</h5>
        {
          source.missions && source.missions.length ?
            <ul>
              {
                source.missions.map(item => <Item key={item.id} {...item} />)
              }
            </ul> :
            <p className="empty">暂无任务</p>
        }
        <IconAdd.A
          className="append"
          onClick={() => {
            props.onAddMissionClick && props.onAddMissionClick(props.source)
          }}
        />
      </div>
    )
  }

  const renderFooter = () => (
    <footer>
      <span>时间周期</span>
      {new Date(source.createTime).format('yyyy年MM月dd日 hh:mm')}
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