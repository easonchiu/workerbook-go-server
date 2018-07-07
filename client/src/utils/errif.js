import Toast from 'src/components/toast'

class Err {
  static errMsg = ''

  static Handle() {
    if (this.errMsg !== '') {
      Toast.error(this.errMsg)
      this.errMsg = ''
      return true
    }
    return false
  }

  static IfEmpty(str, msg) {
    if (this.errMsg === '' && str.trim() === '') {
      this.errMsg = msg
    }
  }

  static IfLenMoreThen(str, len, msg) {
    if (this.errMsg === '' && str.trim().length > len) {
      this.errMsg = msg
    }
  }

  static IfLenLessThen(str, len, msg) {
    if (this.errMsg === '' && str.trim().length < len) {
      this.errMsg = msg
    }
  }
}

export default Err