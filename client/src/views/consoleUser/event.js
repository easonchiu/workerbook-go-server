import Toast from 'src/components/toast'
import Loading from 'src/components/loading'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      Loading.show()
      await Promise.all([
        this.props.$user.fetchList({
          skip: pager * 30 - 30,
          limit: 30,
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

  // 新增人员提交
  onFormSubmit = async data => {
    try {
      await this.props.$user.create(data)
      this.onCloseDialog()
      await this.fetchData()
      Toast.success('添加成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
  }

  // 修改人员提交
  onFormEditSubmit = async data => {
    try {
      await this.props.$user.update(data)
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
    this.userDialog && this.userDialog.$clear()
    this.onOpenDialog()
  }

  // 编辑按钮点击
  onEditClick = async data => {
    try {
      Loading.show()
      const res = await this.props.$user.fetchOneById(data.id)
      this.userDialog && this.userDialog.$fill(res)
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
      userDialogVisible: true
    })
  }

  // 关闭弹层
  onCloseDialog = () => {
    this.setState({
      userDialogVisible: false
    })
  }

}