import Toast from 'src/components/toast'
import Loading from 'src/components/loading'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      await Promise.all([
        this.props.$department.fetchList({
          skip: pager * 30 - 30,
          limit: 30,
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

  // 新增部门提交
  onFormSubmit = async data => {
    try {
      await this.props.$department.create(data)
      this.onCloseDialog()
      await this.fetchData()
      Toast.success('添加成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
  }

  // 修改部门提交
  onFormEditSubmit = async data => {
    try {
      await this.props.$department.update(data)
      this.onCloseDialog()
      await this.fetchData()
      Toast.success('修改成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
  }

  // 添加按钮点击
  onAppendClick = () => {
    this.departmentDialog && this.departmentDialog.$clear()
    this.onOpenDialog()
  }

  // 编辑按钮点击
  onEditClick = async data => {
    try {
      const res = await this.props.$department.fetchOneById(data.id)
      this.departmentDialog && this.departmentDialog.$fill(res)
      this.onOpenDialog()
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
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
      departmentDialogVisible: false
    })
  }

}