import Err from 'src/utils/errif'

export default class Event {

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