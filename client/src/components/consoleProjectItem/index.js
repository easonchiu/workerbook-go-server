import './style'
import React from 'react'
import classNames from 'classnames'

import IconSetting from 'src/components/svg/setting'
import IconAdd from 'src/components/svg/add'

const ProjectItem = props => {
  const source = props.source || {}
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
      <div className="missions">
        <h6>包含任务</h6>
        <ul>
          <li>
            <p>前端页面开发</p>
            <a href="javascript:;">编辑</a>
          </li>
          <li>
            <p>接口设计</p>
            <a href="javascript:;">编辑</a>
          </li>
          <li>
            <p>测试及上线</p>
            <a href="javascript:;">编辑</a>
          </li>
        </ul>
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
      <div className={classNames('description', { empty: !source.description })}>
        <p>{source.description || '暂无说明内容'}</p>
      </div>
      <div className="tools">
        <a
          href="javascript:;"
          className="edit"
          onClick={() => {
            props.onEditClick && props.onEditClick(props.source)
          }}
        >
          <IconSetting />
        </a>
      </div>
      <footer>
        <span>时间周期</span>
        {new Date(source.deadline).format('yyyy年MM月dd日 hh:mm')}
        {' ~ '}
        {new Date(source.deadline).format('MM月dd日 hh:mm')}
      </footer>
    </div>
  )
}
export default ProjectItem