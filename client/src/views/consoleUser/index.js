import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import Input from 'src/components/input'
import Select from 'src/components/select'
import Form from 'src/containers/form'
import Pager from 'src/components/pager'
import MainDialog from 'src/containers/mainDialog'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleUser extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      userDialogVisible: false,
      userId: '',
      nickname: '',
      username: '',
      departmentId: '',
      title: '',
      role: '',
      status: '',
      password: '',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
    this.evt.fetchDepartments()
  }

  renderDialog() {
    const { select } = this.props.department$
    return (
      <MainDialog
        className="dialog-console-edit-user"
        title={this.state.userId ? '修改人员' : '添加人员'}
        visible={this.state.userDialogVisible}
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
            <Select
              value={this.state.departmentId}
              onClick={this.evt.onFormDepartmentChange}
            >
              {
                select.list.map(item => (
                  <Select.Option key={item.id} value={item.id}>
                    {item.name}
                  </Select.Option>
                ))
              }
            </Select>
          </Form.Row>

          <Form.Row label="职称">
            <Input
              name="title"
              value={this.state.title}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>

          <Form.Row label="职位">
            <Select
              value={this.state.role}
              onClick={this.evt.onFormRoleChange}
            >
              <Select.Option value={1}>开发者</Select.Option>
              <Select.Option value={2}>部门管理者</Select.Option>
              <Select.Option value={3}>观察者</Select.Option>
            </Select>
          </Form.Row>

          {
            this.state.userId ?
              <Form.Row label="状态">
                <Input
                  name="role"
                  value={this.state.status}
                  onChange={this.evt.onFormChange}
                />
              </Form.Row> :
              null
          }

          {
            !this.state.userId ?
              <Form.Row label="帐号">
                <Input
                  name="username"
                  value={this.state.username}
                  onChange={this.evt.onFormChange}
                />
              </Form.Row> :
              null
          }

          {
            !this.state.userId ?
              <Form.Row label="初始密码">
                <Input
                  name="password"
                  value={this.state.password}
                  onChange={this.evt.onFormChange}
                />
              </Form.Row> :
              null
          }

          <Form.Row>
            {
              this.state.userId ?
                <Button onClick={this.evt.onFormEditSubmit}>
                  修改
                </Button> :
                <Button onClick={this.evt.onFormSubmit}>
                  提交
                </Button>
            }
          </Form.Row>
        </Form>
      </MainDialog>
    )
  }

  render() {
    const { users } = this.props.user$
    const header = (
      <tr>
        <td>编号</td>
        <td>姓名</td>
        <td>帐号</td>
        <td>部门</td>
        <td>职称</td>
        <td>职位</td>
        <td>状态</td>
        <td>创建时间</td>
        <td>操作</td>
      </tr>
    )
    const roles = {
      1: '开发者',
      2: '部门管理者',
      3: '观察者',
    }
    const body = users.list.map((res, i) => (
      <tr key={res.id}>
        <td>{users.skip + i + 1}</td>
        <td>{res.nickname}</td>
        <td>{res.username}</td>
        <td>{res.departmentName}</td>
        <td>{res.title}</td>
        <td>{roles[res.role]}</td>
        <td>{res.status}</td>
        <td>{new Date(res.createTime).format('yyyy-MM-dd hh:mm')}</td>
        <td className="c">
          <a href="javascript:;" onClick={() => this.evt.onEditClick(res)}>
            编辑
          </a>
        </td>
      </tr>
    ))
    return (
      <div className="console-user">
        <header>
          <h1>人员管理</h1>
          <Button onClick={this.evt.onAppendClick}>添加</Button>
        </header>
        <table className="console-table">
          <thead>{header}</thead>
          <tbody>{body}</tbody>
        </table>
        <Pager
          current={users.skip / users.limit + 1}
          max={Math.ceil(users.count / users.limit)}
          onClick={this.evt.onPageClick}
        />
        {this.renderDialog()}
      </div>
    )
  }
}

export default ConsoleUser