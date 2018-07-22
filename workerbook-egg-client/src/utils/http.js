import axios from 'axios'
import { clearToken, getToken } from 'src/utils/token'

// type an error
function HttpError(message, data) {
  this.message = message
  this.name = 'HttpError'
  this.data = data || null
}
HttpError.prototype = new Error()
HttpError.prototype.constructor = HttpError

const config = {
  production: '/',
  development: '/proxy',
  test: '/'
}

/**
 * 获取config配置中的请求前置路径
 */
const baseURL = config[process.env.PACKAGE] ? config['development'] : config[process.env.PACKAGE]

/**
 * 配置axios
 */
const http = axios.create({
  baseURL,
  headers: {
    Accept: 'application/json;version=3.0;compress=false',
    'Content-Type': 'application/json;charset=utf-8'
  },
  data: {}
})

/**
 * 请求拦截器，在发起请求之前
 */
http.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    config.headers.authorization = 'Bearer ' + token
  }
  return config
})

/**
 * 接口响应拦截器，在接口响应之后
 */
http.interceptors.response.use(
  config => {
    // success handle
    if (config.status === 204 || config.data.resCode === '000000') {
      return config.data.data
    }
    // need to login, token is overdue or empty
    else if (config.status === 401) {
      clearToken()
      return false
    }
    // return reject error
    return Promise.reject(new HttpError(config.data.resMsg, config.resCode))
  },
  error => {
    return Promise.reject(new HttpError(error.response.data.resMsg || '系统错误', error.response.data.resCode))
  }
)

export default http