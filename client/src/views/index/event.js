
export default class Event {

  fetchData = async () => {
    await Promise.all([
      this.props.$daily.fetchDailyListByDay(),
      this.props.$group.fetchList(),
    ])
  }

}
