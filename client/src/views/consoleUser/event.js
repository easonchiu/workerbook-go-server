import Toast from 'src/components/toast'
import Err from 'src/utils/errif'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      await Promise.all([
        this.props.$user.fetchList({
          skip: pager * 10 - 10,
          limit: 10,
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
      departmentId: e
    })
  }

  // 表单职位字段修改
  onFormRoleChange = e => {
    this.setState({
      role: e
    })
  }

  // 新增人员提交
  onFormSubmit = async () => {
    Err.IfEmpty(this.state.nickname, '姓名不能为空')
    Err.IfEmpty(this.state.departmentId, '请选择部门')
    Err.IfEmpty(this.state.title, '职称不能为空')
    Err.IfEmpty(this.state.role, '请选择职位')
    Err.IfEmpty(this.state.username, '登录帐号不能为空')
    Err.IfEmpty(this.state.password, '初始密码不能为空')

    if (!Err.Handle()) {
      try {
        await this.props.$user.create({
          nickname: this.state.nickname,
          departmentId: this.state.departmentId,
          title: this.state.title,
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
    Err.IfEmpty(this.state.title, '职称不能为空')
    Err.IfEmpty(this.state.role, '请选择职位')

    if (!Err.Handle()) {
      try {
        await this.props.$user.update({
          id: this.state.userId,
          nickname: this.state.nickname,
          departmentId: this.state.departmentId,
          title: this.state.title,
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
      title: data.title,
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
      departmentId: '',
      title: '',
      role: '',
      status: '',
      password: '',
    })
  }

}