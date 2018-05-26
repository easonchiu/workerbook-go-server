export default class Event {

  // 获取首页需要的数据
  fetchData = async () => {
    try {
      await Promise.all([
        this.props.$user.myProfile(),
        this.props.$daily.mine(),
        this.props.$daily.fetchListByDay(),
        this.props.$group.fetchList(),
        this.props.$user.fetchList(),
        this.props.$project.fetchList({
          status: 1
        }),
      ])
    }
    catch (err) {
      alert(err.message)
    }
  }

  // 侧栏分组点击
  groupClick = gid => {
    this.props.$user.fetchList({ gid })
  }

  // 日报内容编辑
  dailyWriterChange = e => {
    this.setState({
      record: e.target.value
    })
  }

  // 重新编辑我今天写的日报
  editDailyClick = id => {

  }

  // 删除我今天写的日报
  deleteDailyClick = async id => {
    await this.props.$user.deleteDailyItem({ id })
  }

  // 进度编辑
  progressChange = e => {
    this.setState({
      progress: e.target.value
    })
  }

  // 项目归属编辑
  dailyProjectChange = e => {
    this.setState({
      project: e.target.value
    })
  }

  // 添加日报（发布按钮点击）
  appendDaily = async () => {
    try {
      await this.props.$user.appendDailyItem({
        record: this.state.record,
        progress: this.state.progress,
        project: this.state.project,
      })
    }
    catch (err) {
      alert(err.message)
    }
  }

}