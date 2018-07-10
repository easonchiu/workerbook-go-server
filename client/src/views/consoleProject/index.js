import './style'
import React from 'react'
import VIEW from 'src/hoc/view'
import ComponentEvent from 'src/hoc/componentEvent'
import Event from './event'

import Button from 'src/components/button'
import ProjectItem from 'src/components/projectItem'
import Input from 'src/components/input'
import Form from 'src/containers/form'
import Pager from 'src/components/pager'
import MainDialog from 'src/containers/mainDialog'
import Select from 'src/components/select'

@VIEW
@ComponentEvent('evt', Event)
class ConsoleProject extends React.PureComponent {
  constructor(props) {
    super(props)
    this.state = {
      projectDialogVisible: false,
      projectId: '',
      name: '',
      deadline: '',
      departments: [],
      description: '',
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
        className="dialog-console-edit-project"
        title={this.state.projectId ? '修改项目' : '添加项目'}
        visible={this.state.projectDialogVisible}
        onClose={this.evt.onCloseDialog}
      >
        <Form>
          <Form.Row label="项目名称">
            <Input
              name="name"
              value={this.state.name}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>

          <Form.Row label="截至时间">
            <Input
              name="deadline"
              value={this.state.deadline}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>

          <Form.Row label="参与部门">
            <Select
              value={this.state.departments}
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

          <Form.Row label="说明">
            <Input
              name="description"
              value={this.state.description}
              onChange={this.evt.onFormChange}
            />
          </Form.Row>

          <Form.Row>
            {
              this.state.projectId ?
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
    const { projects } = this.props.project$
    const row = []
    for (let i = 0; i < projects.list.length; i += 3) {
      row.push(
        <div className="row" key={i}>
          {
            [0, 1, 2].map(j => {
              const item = projects.list[i + j]
              if (item) {
                return <ProjectItem key={j} source={item} />
              }
              return <div key={j} />
            })
          }
        </div>
      )
    }
    return (
      <div className="console-project">
        <header>
          <h1>项目管理</h1>
          <Button onClick={this.evt.onAppendClick}>添加</Button>
        </header>
        <div className="list">
          {row}
        </div>
        <Pager
          current={projects.skip / projects.limit + 1}
          max={Math.ceil(projects.count / projects.limit)}
          onClick={this.evt.onPageClick}
        />
        {this.renderDialog()}
      </div>
    )
  }
}

export default ConsoleProject