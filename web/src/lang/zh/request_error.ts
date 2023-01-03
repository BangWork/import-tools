const requestError = {
  networkError: '后端接口连接异常',
  timeout: '系统接口请求超时',
  requestFailed: '接口异常',
  default: '未知错误',
  400: '错误的请求',
  401: '未授权，请重新登录',
  403: '拒绝访问',
  404: '请求错误,未找到该资源',
  405: '请求方法未允许',
  408: '请求超时',
  500: '服务器端出错',
  501: '网络未实现',
  502: '网络错误',
  503: '服务不可用',
  504: '网络超时',
  505: 'http版本不支持该请求',
  other: '其他连接错误 --{{code}}',
};

export default requestError;
