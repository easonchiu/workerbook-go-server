import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import Input from 'src/components/input'
import Form from 'src/containers/form'
import Pager from 'src/components/pager'
import MainDialog from 'src/containers/mainDialog'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleUserList extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      pager: 1,
      appendDialogVisible: false,
      nickname: '赵志达',
      username: 'eason',
      departmentId: '5b424feeaea6f431c2655006',
      role: '1',
      password: '123456',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  renderDialog() {
    return (
      <MainDialog
        className="dialog-console-edit-user"
        title="添加人员"
        visible={this.state.appendDialogVisible}
        onClose={this.evt.onCloseDialog}
      >
        <Form>
          <Form.Row label="姓名">
            <Input
              name="nickname"
              value={this.state.nickname}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>
          <Form.Row label="部门">
            <Input
              name="departmentId"
              value={this.state.departmentId}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>
          <Form.Row label="部门主管">
            <Input
              name="role"
              value={this.state.role}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>
          <Form.Row label="帐号">
            <Input
              name="username"
              value={this.state.username}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>
          <Form.Row label="初始密码">
            <Input
              name="password"
              value={this.state.password}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>
          <Form.Row>
            <Button onClick={this.evt.onFormSubmit}>
              提交
            </Button>
          </Form.Row>
        </Form>
      </MainDialog>
    )
  }

  render() {
    const { users } = this.props.console$
    const header = (
      <tr>
        <td>编号</td>
        <td>姓名</td>
        <td>帐号</td>
        <td>部门</td>
        <td>状态</td>
        <td>操作</td>
      </tr>
    )
    const body = users.list.map((res, i) => (
      <tr key={res.id}>
        <td>{i}</td>
        <td>{res.nickname}</td>
        <td>{res.username}</td>
        <td>{res.departmentId}</td>
        <td>{res.role}</td>
        <td className="c"><a href="javascript:;">编辑</a></td>
      </tr>
    ))
    return (
      <div className="console-user-list">
        <header>
          <h1>人员管理</h1>
          <Button onClick={this.evt.onOpenDialog}>添加</Button>
        </header>
        <table className="console-table">
          <thead>{header}</thead>
          <tbody>{body}</tbody>
        </table>
        <Pager
          current={this.state.pager}
          max={12}
          onClick={this.evt.onPageClick}
        />
        {this.renderDialog()}
      </div>
    )
  }
}

export default ConsoleUserList