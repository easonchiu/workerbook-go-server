import './style'
import React, { PureComponent } from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Wrapper from 'src/containers/wrapper'
import Dialog from 'src/containers/dialog'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  componentDidMount() {
    this.evt.fetchData()
  }

  render(props, state) {
    const profile = this.props.user$.profile

    return (
      <div className="view-chart">
        <Wrapper.Header nav="chart" profile={profile} />

        <Wrapper.Body>
          xxx
        </Wrapper.Body>

        <Wrapper.Footer />

        <Dialog />
      </div>
    )
  }
}
