import Toast from 'src/components/toast'
import Err from 'src/utils/errif'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      await Promise.all([
        this.props.$project.fetchList({
          skip: pager * 9 - 9,
          limit: 9
        })
      ])
    }
    catch (err) {
      Toast.error(err.message)
    }
  }

  // 获取所有部门信息
  fetchDepartments = async () => {
    try {
      await this.props.$department.fetchSelectList()
    }
    catch (err) {
      Toast.error(err.message)
    }
  }

  // 翻页
  onPageClick = pager => {
    this.fetchData(pager)
  }

  // 表单字段修改
  onFormChange = e => {
    const key = e.target.name
    this.setState({
      [key]: e.target.value
    })
  }

  // 表单部门字段修改
  onFormDepartmentChange = e => {
    this.setState({
      departments: [e]
    })
  }

  // 新增项目提交
  onFormSubmit = async () => {
    this.state.deadline = new Date(2018, 11, 20)
    Err.IfEmpty(this.state.name, '项目名称不能为空')
    Err.IfEmpty(this.state.deadline, '截至时间不能为空')
    Err.IfEmptyArr(this.state.departments, '参与部门不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$project.create({
          name: this.state.name,
          deadline: this.state.deadline,
          departments: this.state.departments,
          description: this.state.description,
        })
        this.onCloseDialog()
        await this.fetchData()
        Toast.success('添加成功')
      }
      catch (err) {
        Toast.error(err.message)
      }
    }
  }

  // 修改项目提交
  onFormEditSubmit = async () => {
    Err.IfEmpty(this.state.name, '项目名称不能为空')
    Err.IfLenMoreThen(this.state.name, 30, '项目名称不能超过30个字节')
    Err.IfEmpty(this.state.deadline, '截至时间不能为空')
    Err.IfEmptyArr(this.state.departments, '参与部门不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$project.update({
          id: this.state.projectId,
          name: this.state.name,
          deadline: this.state.deadline,
          departments: this.state.departments,
          description: this.state.description,
        })
        this.onCloseDialog()
        await this.fetchData()
        Toast.success('修改成功')
      }
      catch (err) {
        Toast.error(err.message)
      }
    }
  }

  // 添加按钮点击
  onAppendClick = () => {
    this.setState({
      projectId: ''
    })
    this.onOpenDialog()
  }

  // 编辑按钮点击
  onEditClick = data => {
    this.setState({
      name: data.name,
      projectId: data.id,
    })
    this.onOpenDialog()
  }

  // 打开弹层
  onOpenDialog = () => {
    this.setState({
      projectDialogVisible: true
    })
  }

  // 关闭弹层
  onCloseDialog = () => {
    this.setState({
      projectDialogVisible: false,
      name: '',
    })
  }

}