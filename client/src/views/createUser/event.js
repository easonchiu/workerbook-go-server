
export default class Event {

  fetchData = async () => {
    await this.props.$group.fetchList()
  }

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
      await this.props.$user.create({
        username: this.state.username,
        password: this.state.password,
        nickname: this.state.nickname,
        gid: this.state.gid,
        role: this.state.role,
        email: this.state.email,
        mobile: this.state.mobile,
      })
    }
    catch (err) {
      alert(err.message)
    }

  }

}
