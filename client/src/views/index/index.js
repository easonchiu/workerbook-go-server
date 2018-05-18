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

    }
  }

  componentDidMount() {
    this.evt.fetchData()
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
        {
          list.map(item => <li key={item.id}>{item.name}-{item.count}</li>)
        }
      </ul>
    )
  }

  render(props, state) {
    return (
      <main className={'view-index'}>

        <h2>Daily list.</h2>
        {this.renderDailyList()}

        <h2>Group list.</h2>
        {this.renderGroupList()}

      </main>
    )
  }
}
