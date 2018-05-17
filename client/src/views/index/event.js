
export default class Event {

  fetchData = async () => {
    await this.props.$daily.fetchUsersDailyList()
  }

}
