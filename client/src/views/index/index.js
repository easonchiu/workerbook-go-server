import './style'
import React, { PureComponent } from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Wrapper from 'src/containers/wrapper'
import DailyList from 'src/containers/dailyList'

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
            </li>
          ))
        }
      </ul>
    )
  }

  renderDailyWriter() {
    const list = this.props.project$.list

    return (
      <div>
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
    const list = this.props.daily$.list
    return <DailyList list={list} />
  }

  renderGroupList() {
    const list = this.props.group$.list
    return (
      <ul>
        <li>
          <a href="javascript:;" onClick={this.evt.groupClick.bind(null, '')}>
            全部
          </a>
        </li>
        {
          list.map(item => (
            <li key={item.id}>
              <a href="javascript:;" onClick={this.evt.groupClick.bind(null, item.id)}>
                {item.name}-{item.count}
              </a>
            </li>
          ))
        }
      </ul>
    )
  }

  renderUserList() {
    const list = this.props.user$.list

    return (
      <ul>
        {
          list.map(item => (
            <li key={item.id}>
              <a href="javascript:;">
                {item.nickname} ({item.groupName})
              </a>
            </li>
          ))
        }
      </ul>
    )
  }

  renderProjectList() {
    const list = this.props.project$.list

    return (
      <ul>
        {
          list.map(item => (
            <li key={item.id}>
              <a href="javascript:;">
                {item.name}
              </a>
            </li>
          ))
        }
      </ul>
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

            <h2>Write daily.</h2>
            {this.renderDailyWriter()}

            <h2>Daily list.</h2>
            {this.renderDailyList()}
          </main>

          <aside className="app-body__aside">
            <h2>Group list.</h2>
            {this.renderGroupList()}

            <h2>User list.</h2>
            {this.renderUserList()}

            <h2>Project list.</h2>
            {this.renderProjectList()}
          </aside>
        </div>

      </div>
    )
  }
}
