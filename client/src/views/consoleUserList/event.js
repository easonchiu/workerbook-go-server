import Err from 'src/utils/errif'

export default class Event {
  fetchData = async () => {
    try {
      await Promise.all([
        this.props.$console.fetchUsersList({
          skip: 0,
          limit: 0,
          isConsole: true,
        })
      ])
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

  onPageClick = p => {
    this.setState({
      pager: p
    })
  }

  onFormChange = e => {
    const key = e.target.name
    this.setState({
      [key]: e.target.value
    })
  }

  onFormSubmit = () => {
    Err.IfEmpty(this.state.nickname, '姓名不能为空')
    Err.IfEmpty(this.state.departmentId, '请选择部门')
    Err.IfEmpty(this.state.username, '登录帐号不能为空')
    Err.IfEmpty(this.state.password, '初始密码不能为空')

    if (!Err.Handle()) {
      this.props.$user.create({
        nickname: this.state.nickname,
        departmentId: this.state.departmentId,
        role: this.state.role,
        username: this.state.username,
        password: this.state.password,
      })
    }

  }

  onOpenDialog = () => {
    this.setState({
      appendDialogVisible: true
    })
  }

  onCloseDialog = () => {
    this.setState({
      appendDialogVisible: false
    })
  }

}