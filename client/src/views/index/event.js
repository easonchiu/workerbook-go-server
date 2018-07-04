export default class Event {

  // 获取首页需要的数据
  fetchData = async () => {
    try {
      await Promise.all([
        this.props.$user.fetchMyProfile(),
        this.props.$user.fetchMyTodayDaily(),
        this.props.$daily.fetchListByDay(),
        this.props.$group.fetchList(),
        this.props.$user.fetchList(),
        this.props.$project.fetchList({
          status: 1
        }),
      ])
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

  // 更新数据
  reload = async () => {
    try {
      await Promise.all([
        this.props.$user.fetchMyTodayDaily(),
        this.props.$daily.fetchListByDay(),
      ])
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

  // 侧栏分组翻页
  groupPageChange = async page => {
    try {
      const { limit } = this.props.group$.list
      await this.props.$group.fetchList({ skip: (page - 1) * limit })
    }
    catch (err) {
      alert(err.message)
    }
  }

  // 侧栏分组点击
  groupClick = async gid => {
    try {
      await this.props.$user.fetchList({ gid })
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

  // 侧栏项目点击
  projectClick = async pid => {
    console.log(pid)
  }

  // 删除我今天写的日报
  deleteDailyClick = async id => {
    try {
      await this.props.$user.deleteDailyItem({ id })
      this.reload()
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

  // 添加日报（发布按钮点击）
  appendDaily = async ({ record, progress, project }) => {
    try {
      await this.props.$user.appendDailyItem({
        record,
        progress,
        project,
      })
      this.$myDailyWriter && this.$myDailyWriter.clear()
      this.reload()
    }
    catch (err) {
      console.error(err)
      alert(err.message)
    }
  }

}