import './style'
import React, { PureComponent } from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Wrapper from 'src/containers/wrapper'
import MainDailyList from 'src/containers/mainDailyList'
import AsideGroupList from 'src/containers/asideGroupList'
import AsideUserList from 'src/containers/asideUserList'
import AsideProjectList from 'src/containers/asideProjectList'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  constructor(props) {
    super(props)

    this.state = {
      record: '',
      progress: 0,
      project: '',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  renderMyDaily() {
    const list = this.props.daily$.mine
    return (
      <ul>
        {
          list.map(i => (
            <li key={i.id}>
              {i.progress}%{' - '}
              {
                i.pname ? i.pname + ' - ' : null
              }
              {i.record}
              <a href="javascript:;" onClick={this.evt.editDailyClick.bind(null, i.id)}>
                编辑
              </a>
              <a href="javascript:;" onClick={this.evt.deleteDailyClick.bind(null, i.id)}>
                删除
              </a>
            </li>
          ))
        }
      </ul>
    )
  }

  renderDailyWriter() {
    const list = this.props.project$.list

    return (
      <div className="daily-writer">
        <input
          type="text"
          placeholder="writer"
          value={this.state.record}
          onChange={this.evt.dailyWriterChange}
        />

        <input
          type="text"
          placeholder="progress"
          value={this.state.progress}
          onChange={this.evt.progressChange}
        />

        <select
          value={this.state.project}
          onChange={this.evt.dailyProjectChange}
        >
          <option value={0}>请选择</option>
          {
            list.map(item => (
              <option
                key={item.id}
                value={item.id}
              >
                {item.name}
              </option>
            ))
          }
        </select>

        <button onClick={this.evt.appendDaily}>
          发布
        </button>
      </div>
    )
  }

  renderDailyList() {
    return <MainDailyList list={this.props.daily$.list} />
  }

  renderGroupList() {
    return (
      <AsideGroupList
        list={this.props.group$.list}
        active={''}
        itemClick={this.evt.groupClick}
      />
    )
  }

  renderUserList() {
    return (
      <AsideUserList
        list={this.props.user$.list}
        itemClick={this.evt.groupClick}
      />
    )
  }

  renderProjectList() {
    return (
      <AsideProjectList
        list={this.props.project$.list}
        active={''}
        itemClick={this.evt.groupClick}
      />
    )
  }

  render(props, state) {
    const profile = this.props.user$.profile

    return (
      <div className="view-index">
        <Wrapper.Header profile={profile} />

        <div className="app-body">
          <main className="app-body__main">
            {this.renderMyDaily()}

            {this.renderDailyWriter()}

            {this.renderDailyList()}
          </main>

          <aside className="app-body__aside">
            {this.renderGroupList()}
            {this.renderUserList()}
            {this.renderProjectList()}
          </aside>
        </div>

      </div>
    )
  }
}
