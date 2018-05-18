
export default class Event {

  fetchData = async () => {
    try {
      await Promise.all([
        this.props.$user.myProfile(),
        this.props.$daily.fetchDailyListByDay(),
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

  groupClick = gid => {
    this.props.$user.fetchList({
      gid
    })
  }

  dailyWriterChange = e => {
    this.setState({
      record: e.target.value
    })
  }

  dailyProjectChange = e => {
    this.setState({
      project: e.target.value
    })
  }

  appendDaily = async () => {
    try {
      await this.props.$daily.append({
        record: this.state.record,
        project: this.state.project,
      })
    }
    catch (err) {
      alert(err.message)
    }
  }

}