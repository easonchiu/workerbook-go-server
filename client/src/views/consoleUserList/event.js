import Toast from 'src/components/toast'
import Err from 'src/utils/errif'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      await Promise.all([
        this.props.$console.fetchUsersList({
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

  // 新增人员提交
  onFormSubmit = async () => {
    Err.IfEmpty(this.state.nickname, '姓名不能为空')
    Err.IfEmpty(this.state.departmentId, '请选择部门')
    Err.IfEmpty(this.state.username, '登录帐号不能为空')
    Err.IfEmpty(this.state.password, '初始密码不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$console.createUser({
          nickname: this.state.nickname,
          departmentId: this.state.departmentId,
          role: this.state.role,
          username: this.state.username,
          password: this.state.password,
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

  // 修改人员提交
  onFormEditSubmit = async () => {
    Err.IfEmpty(this.state.nickname, '姓名不能为空')
    Err.IfEmpty(this.state.departmentId, '请选择部门')

    if (!Err.Handle()) {
      try {
        await this.props.$console.updateUser({
          id: this.state.userId,
          nickname: this.state.nickname,
          departmentId: this.state.departmentId,
          role: this.state.role,
          status: this.state.status,
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
      userId: ''
    })
    this.onOpenDialog()
  }

  // 编辑按钮点击
  onEditClick = data => {
    this.setState({
      nickname: data.nickname,
      username: data.username,
      departmentId: data.departmentId,
      role: data.role,
      status: data.status,
      userId: data.id,
    })
    this.onOpenDialog()
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
      userDialogVisible: false,
      nickname: '',
      username: '',
      departmentId: '5b424feeaea6f431c2655006',
      role: '',
      status: '',
      password: '',
    })
  }

}