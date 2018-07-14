import Toast from 'src/components/toast'
import Loading from 'src/components/loading'

export default class Event {
  fetchData = async (pager = 1) => {
    try {
      Loading.show()
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
    finally {
      Loading.hide()
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

  // 新增项目提交
  onProjectFormSubmit = async data => {
    data.deadline = new Date(2018, 11, 20)
    try {
      Loading.show()
      await this.props.$project.create(data)
      this.onCloseProjectDialog()
      await this.fetchData()
      Toast.success('添加成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

  // 修改项目提交
  onProjectFormEditSubmit = async data => {
    data.deadline = new Date(2018, 11, 20)
    try {
      Loading.show()
      await this.props.$project.update(data)
      this.onCloseProjectDialog()
      await this.fetchData()
      Toast.success('修改成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

  // 项目添加按钮点击
  onAppendProjectClick = () => {
    this.projectDialog && this.projectDialog.$clear()
    this.onOpenProjectDialog()
  }

  // 项目编辑按钮点击
  onProjectEditClick = async data => {
    try {
      Loading.show()
      const res = await this.props.$project.fetchOneById(data.id)
      res.departments = res.departments ? res.departments.map(i => i.id) : []
      this.projectDialog && this.projectDialog.$fill(res)
      this.onOpenProjectDialog()
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

  // 项目删除按钮点击
  onProjectDeleteClick = data => {
    if (data.name && data.id) {
      this.setState({
        projectDelDialogVisible: true,
        projectDelDialogData: data,
      })
    }
    else {
      Toast.error('系统错误')
    }
  }

  // 关闭项目删除弹层
  onCloseDelProjectDialog = () => {
    this.setState({
      projectDelDialogVisible: false
    })
  }

  // 确定删除项目
  onDelProject = async data => {
    if (data.name && data.id) {
      try {
        Loading.show()
        await this.props.$project.del(data.id)
        Toast.success('删除成功')
      }
      catch (err) {
        Toast.error(err.message)
      }
      finally {
        Loading.hide()
      }
    }
    else {
      Toast.error('系统错误')
    }
  }

  // 打开项目弹层
  onOpenProjectDialog = () => {
    this.setState({
      projectDialogVisible: true
    })
  }

  // 关闭项目弹层
  onCloseProjectDialog = () => {
    this.setState({
      projectDialogVisible: false
    })
  }

  // ------------任务-------------

  // 添加任务按钮点击
  onAppendMissionClick = project => {
    if (this.missionDialog) {
      this.missionDialog.$clear()
      this.missionDialog.$projectId(project.id)
    }
    this.onOpenMissionDialog()
  }

  // 编辑任务点击
  onMissionEditClick = async (data, project) => {
    try {
      Loading.show()
      const res = await this.props.$mission.fetchOneById(data.id)
      if (this.missionDialog) {
        this.missionDialog.$fill(res)
        this.missionDialog.$projectId(project.id)
      }
      this.onOpenMissionDialog()
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

  // 打开任务弹层
  onOpenMissionDialog = () => {
    this.setState({
      missionDialogVisible: true
    })
  }

  // 关闭任务弹层
  onCloseMissionDialog = () => {
    this.setState({
      missionDialogVisible: false
    })
  }

  // 新增任务提交
  onMissionFormSubmit = async data => {
    data.deadline = new Date(2018, 11, 20)
    try {
      Loading.show()
      await this.props.$mission.create(data)
      this.onCloseMissionDialog()
      await this.fetchData()
      Toast.success('添加成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

  // 修改任务提交
  onMissionFormEditSubmit = async data => {
    data.deadline = new Date(2018, 11, 20)
    try {
      Loading.show()
      await this.props.$mission.update(data)
      this.onCloseMissionDialog()
      await this.fetchData()
      Toast.success('修改成功')
    }
    catch (err) {
      Toast.error(err.message)
    }
    finally {
      Loading.hide()
    }
  }

}