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
      username: '',
      password: '',
      nickname: '',
      gid: '',
      role: 1,
      email: '',
      mobile: '',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  render(props, state) {
    const groups = props.group$.list

    return (
      <main className="view-create-user">
        <form>
          <input
            type="text"
            placeholder="用户名"
            name="username"
            value={this.state.username}
            onChange={this.evt.formValueChange}
          />

          <input
            type="text"
            placeholder="密码"
            name="password"
            value={this.state.password}
            onChange={this.evt.formValueChange}
          />

          <input
            type="text"
            placeholder="昵称"
            name="nickname"
            value={this.state.nickname}
            onChange={this.evt.formValueChange}
          />

          <select
            value={this.state.gid}
            name="gid"
            onChange={this.evt.formValueChange}
          >
            <option value={0}>请选择</option>
            {
              groups.map(item => {
                return (
                  <option key={item.id} value={item.id}>
                    {item.name}
                  </option>
                )
              })
            }
          </select>

          <select
            value={this.state.role}
            name="role"
            onChange={this.evt.formValueChange}
          >
            <option value={1}>普通成员</option>
            <option value={2}>部门Leader</option>
            <option value={99}>管理员</option>
          </select>

          <input
            type="text"
            placeholder="Email"
            name="email"
            value={this.state.email}
            onChange={this.evt.formValueChange}
          />

          <input
            type="text"
            placeholder="手机号"
            name="mobile"
            value={this.state.mobile}
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
