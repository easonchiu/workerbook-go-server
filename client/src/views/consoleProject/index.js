import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import ProjectItem from 'src/components/consoleProjectItem'
import Pager from 'src/components/pager'
import ConsoleProjectDialog from 'src/components/consoleProjectDialog'
import ConsoleMissionDialog from 'src/components/consoleMissionDialog'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleProject extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      projectDialogVisible: false,
      missionDialogVisible: false,
    }
  }

  componentDidMount() {
    this.evt.fetchData()
    this.evt.fetchDepartments()
  }

  renderProjectDialog() {
    const { select } = this.props.department$
    return (
      <ConsoleProjectDialog
        ref={r => { this.projectDialog = r }}
        departments={select ? select.list || [] : []}
        visible={this.state.projectDialogVisible}
        onClose={this.evt.onCloseProjectDialog}
        onSubmit={this.evt.onProjectFormSubmit}
        onEditSubmit={this.evt.onProjectFormEditSubmit}
      />
    )
  }

  renderMissionDialog() {
    return (
      <ConsoleMissionDialog
        ref={r => { this.missionDialog = r }}
        visible={this.state.missionDialogVisible}
        onClose={this.evt.onCloseMissionDialog}
        onSubmit={this.evt.onMissionFormSubmit}
        onEditSubmit={this.evt.onMissionFormEditSubmit}
      />
    )
  }

  render() {
    const { projects } = this.props.project$
    const row = []
    for (let i = 0; i < projects.list.length; i += 3) {
      row.push(
        <div className="row" key={i}>
          {
            [0, 1, 2].map(j => {
              const item = projects.list[i + j]
              if (item) {
                return (
                  <ProjectItem
                    key={j}
                    onAppendMissionClick={this.evt.onAppendMissionClick}
                    onEditClick={this.evt.onProjectEditClick}
                    source={item}
                  />
                )
              }
              return <div className="space" key={j} />
            })
          }
        </div>
      )
    }
    return (
      <div className="console-project">
        <header>
          <h1>项目管理</h1>
          <Button onClick={this.evt.onAppendProjectClick}>添加</Button>
        </header>
        <div className="list">
          {row}
        </div>
        <Pager
          current={projects.skip / projects.limit + 1}
          max={Math.ceil(projects.count / projects.limit)}
          onClick={this.evt.onPageClick}
        />
        {this.renderProjectDialog()}
        {this.renderMissionDialog()}
      </div>
    )
  }
}

export default ConsoleProject