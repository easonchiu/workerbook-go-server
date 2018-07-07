import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import Pager from 'src/components/pager'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleDepartmentList extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      pager: 1
    }
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
      <div className="console-department-list">
        <header>
          <h1>部门管理</h1>
          <Button>添加</Button>
        </header>
        <table className="console-table">
          <thead>{header}</thead>
          <tbody>{body}</tbody>
        </table>
        <Pager current={this.state.pager} max={89} onClick={this.evt.onPageClick} />
      </div>
    )
  }
}

export default ConsoleDepartmentList