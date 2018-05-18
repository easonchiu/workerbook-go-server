
export default class Event {

  formValueChange = e => {
    const name = e.target.name
    if (name) {
      this.setState({
        [name]: e.target.value
      })
    }
  }

  onSubmit = async e => {
    e.preventDefault()

    await this.props.$group.create({
      name: this.state.name,
    })
  }

}
