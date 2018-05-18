import './style'
import React, { PureComponent } from 'react'
import { Link } from 'react-router-dom'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  constructor(props) {
    super(props)

    this.state = {
      record: '',
      project: '',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  renderMyProfile() {
    const profile = this.props.user$.profile

    return (
      <div>
        {profile.nickname} - {profile.groupName}
      </div>
    )
  }

  renderDailyWriter() {
    const list = this.props.project$.list

    return (
      <div>
        <input
          type={'text'}
          placeholder={'writer'}
          value={this.state.record}
          onChange={this.evt.dailyWriterChange}
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

    return (
      <article>
        {
          list.map(item => {
            return (
              <section key={item.id}>
                <h2>uid: {item.uid}</h2>
                <ul>
                  {
                    item.dailyList.map(i => {
                      return (
                        <li key={i.id}>
                          {i.progress}%,
                          {i.pname}
                          {i.record}
                        </li>
                      )
                    })
                  }
                </ul>
              </section>
            )
          })
        }
      </article>
    )
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
    return (
      <main className={'view-index'}>

        <Link to={'/createGroup'}>Create Group</Link>
        {' | '}
        <Link to={'/createUser'}>Create User</Link>
        {' | '}
        <Link to={'/createProject'}>Create Project</Link>

        <h2>My profile.</h2>
        {this.renderMyProfile()}

        <h2>Write daily.</h2>
        {this.renderDailyWriter()}

        <h2>Daily list.</h2>
        {this.renderDailyList()}

        <h2>Group list.</h2>
        {this.renderGroupList()}

        <h2>User list.</h2>
        {this.renderUserList()}

        <h2>Project list.</h2>
        {this.renderProjectList()}

      </main>
    )
  }
}
