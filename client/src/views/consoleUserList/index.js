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
      departmentId: '5a5cd2181ea5a0c406c6a2e8',
      role: '1',
      password: '123456',
    }
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
    const header = (
      <tr>
        <td>编号</td>
        <td>昵称</td>
        <td>用户名</td>
        <td>部门</td>
        <td>状态</td>
        <td>操作</td>
      </tr>
    )
    const body = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(res => (
      <tr key={res}>
        <td>{res}</td>
        <td>王三</td>
        <td>testname</td>
        <td>前端开发</td>
        <td>在职</td>
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