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

  render(props, state) {

    const list = props.daily$.list

    console.log(list)

    return (
      <main className={'view-index'}>

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

      </main>
    )
  }
}
