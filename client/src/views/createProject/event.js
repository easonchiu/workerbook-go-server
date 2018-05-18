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
    try {
      await this.props.$project.create({
        name: this.state.name,
      })
      this.props.history.goBack()
    }
    catch (err) {
      alert(err.message)
    }
  }

}
