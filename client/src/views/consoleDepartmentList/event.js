export default class Event {

  onPageClick = p => {
    this.setState({
      pager: p
    })
  }

}