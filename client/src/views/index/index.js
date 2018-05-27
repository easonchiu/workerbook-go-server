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
import MyDailyWriter from 'src/components/myDailyWriter'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  componentDidMount() {
    this.evt.fetchData()
  }

  // 日报编写区
  renderMyDailyWriter() {
    return (
      <MyDailyWriter
        ref={el => {
          this.$myDailyWriter = el
        }}
        myDailyList={this.props.user$.fetchMyTodayDaily}
        projectList={this.props.project$.list}
        onDeleteItem={this.evt.deleteDailyClick}
        onAppend={this.evt.appendDaily}
      />
    )
  }

  // 主体区的日报列表
  renderDailyList() {
    return <MainDailyList list={this.props.daily$.list} />
  }

  // 侧栏的分组模块
  renderGroupList() {
    return (
      <AsideGroupList
        data={this.props.group$.list}
        active={this.props.user$.activeGroup}
        itemClick={this.evt.groupClick}
      />
    )
  }

  // 侧栏的用户模块
  renderUserList() {
    return (
      <AsideUserList
        list={this.props.user$.list}
        isAll={this.props.user$.activeGroup === ''}
        itemClick={this.evt.groupClick}
      />
    )
  }

  // 侧栏的项目模块
  renderProjectList() {
    return (
      <AsideProjectList
        list={this.props.project$.list}
        active={this.props.daily$.activeProject}
        itemClick={this.evt.projectClick}
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
            {this.renderMyDailyWriter()}
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
