import Toast from 'src/components/toast'
import Err from 'src/utils/errif'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      await Promise.all([
        this.props.$console.fetchDepartmentsList({
          skip: pager * 10 - 10,
          limit: 10,
        })
      ])
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

  // 新增部门提交
  onFormSubmit = async () => {
    Err.IfEmpty(this.state.name, '部门名称不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$console.createDepartment({
          name: this.state.name,
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

  // 修改部门提交
  onFormEditSubmit = async () => {
    Err.IfEmpty(this.state.name, '部门名称不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$console.updateDepartment({
          id: this.state.departmentId,
          name: this.state.name,
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
      departmentId: ''
    })
    this.onOpenDialog()
  }

  // 编辑按钮点击
  onEditClick = data => {
    this.setState({
      name: data.name,
      departmentId: data.id,
    })
    this.onOpenDialog()
  }

  // 打开弹层
  onOpenDialog = () => {
    this.setState({
      departmentDialogVisible: true
    })
  }

  // 关闭弹层
  onCloseDialog = () => {
    this.setState({
      departmentDialogVisible: false,
      name: '',
    })
  }

}