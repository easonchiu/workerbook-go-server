import './style'
import React, { PureComponent } from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

@VIEW
@ComponentEvent('evt', Event)
export default class View extends PureComponent {
  constructor(props) {
    super(props)

    this.state = {
      name: '',
    }
  }

  render(props, state) {
    return (
      <main className={'view-create-group'}>
        <form>
          <input
            type="text"
            placeholder="分组名称"
            name="name"
            value={this.state.name}
            onChange={this.evt.formValueChange}
          />

          <button onClick={this.evt.onSubmit}>
            提交
          </button>
        </form>
      </main>
    )
  }
}
