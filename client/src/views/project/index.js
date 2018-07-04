import './style'
import React, { PureComponent } from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Wrapper from 'src/containers/wrapper'
import ProjectItem from 'src/components/projectItem'
import MissionItem from 'src/components/missionItem'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      showMission: false
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  render(props, state) {
    const profile = this.props.user$.profile

    return (
      <div className="view-project">
        <Wrapper.Header nav="project" profile={profile} />

        <div className="app-body">

          <div className="project-list">
            <header>
              <h1>参与的项目</h1>
            </header>
            <div className="list">
              <ProjectItem onMissionClick={this.evt.click} />
            </div>
          </div>

          <div className="project-list">
            <header>
              <h1>未参与的项目</h1>
            </header>
            <div className="list">
              <ProjectItem />
              <ProjectItem />
              <ProjectItem />
              <ProjectItem />
              <ProjectItem />
              <ProjectItem />
            </div>
          </div>

        </div>

        {
          this.state.showMission ?
            <div className="mission-bar" onClick={this.evt.hide}>
              <div className="inner">
                <header><h2>世界杯活动页面开发</h2></header>
                <MissionItem showJoined />
                <MissionItem showJoined />
                <MissionItem showJoined />
                <MissionItem showJoined />
              </div>
            </div> :
            null
        }

        <Wrapper.Footer />

      </div>
    )
  }
}
