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
class ConsoleDepartment extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      departmentDialogVisible: false,
      departmentId: '',
      name: '',
    }
  }

  componentDidMount() {
    this.evt.fetchData()
  }

  renderDialog() {
    return (
      <MainDialog
        className="dialog-console-edit-department"
        title={this.state.departmentId ? '修改部门' : '添加部门'}
        visible={this.state.departmentDialogVisible}
        onClose={this.evt.onCloseDialog}
      >
        <Form>
          <Form.Row label="部门名称">
            <Input
              name="name"
              value={this.state.name}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>

          <Form.Row>
            {
              this.state.departmentId ?
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
    const { departments } = this.props.department$
    const header = (
      <tr>
        <td>编号</td>
        <td>部门名称</td>
        <td>人数</td>
        <td>创建时间</td>
        <td>操作</td>
      </tr>
    )
    const body = departments.list.map((res, i) => (
      <tr key={res.id}>
        <td>{departments.skip + i + 1}</td>
        <td>{res.name}</td>
        <td>{res.userCount}</td>
        <td>{new Date(res.createTime).format('yyyy-MM-dd hh:mm')}</td>
        <td className="c">
          <a href="javascript:;" onClick={() => this.evt.onEditClick(res)}>
            编辑
          </a>
        </td>
      </tr>
    ))
    return (
      <div className="console-department">
        <header>
          <h1>部门管理</h1>
          <Button onClick={this.evt.onAppendClick}>添加</Button>
        </header>
        <table className="console-table">
          <thead>{header}</thead>
          <tbody>{body}</tbody>
        </table>
        <Pager
          current={departments.skip / departments.limit + 1}
          max={Math.ceil(departments.count / departments.limit)}
          onClick={this.evt.onPageClick}
        />
        {this.renderDialog()}
      </div>
    )
  }
}

export default ConsoleDepartment