import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import Pager from 'src/components/pager'
import ConsoleUserDialog from 'src/components/consoleUserDialog'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleUser extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      userDialogVisible: false
    }
  }

  componentDidMount() {
    this.evt.fetchData()
    this.evt.fetchDepartments()
  }

  renderDialog() {
    const { select } = this.props.department$
    return (
      <ConsoleUserDialog
        ref={r => { this.userDialog = r }}
        departments={select ? select.list || [] : []}
        visible={this.state.userDialogVisible}
        onClose={this.evt.onCloseDialog}
        onSubmit={this.evt.onFormSubmit}
        onEditSubmit={this.evt.onFormEditSubmit}
      />
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
        <td>加入时间</td>
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