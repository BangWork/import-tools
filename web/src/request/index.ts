import axios from 'axios';
import { message } from 'antd';
import i18n, { t } from 'i18next';

import { errorCodeType } from './error_code';

const initConfig = {
  baseURL: import.meta.env.VITE_PROXY_DOMAIN_REAL || window.location.origin,
  timeout: 15000,
  headers: {
    'Accept-Language': i18n.language,
    'Content-Type': 'application/json;charset=utf-8',
  },
};

const requestThen = (config) => {
  // Do you need to set the token in the request header
  // get request mapping params parameter
  if (config.method === 'get' && config.params) {
    let url = config.url + '?';
    for (const propName of Object.keys(config.params)) {
      const value = config.params[propName];
      var part = encodeURIComponent(propName) + '=';
      if (value !== null && typeof value !== 'undefined') {
        if (typeof value === 'object') {
          for (const key of Object.keys(value)) {
            let params = propName + '[' + key + ']';
            var subPart = encodeURIComponent(params) + '=';
            url += subPart + encodeURIComponent(value[key]) + '&';
          }
        } else {
          url += part + encodeURIComponent(value) + '&';
        }
      }
    }
    url = url.slice(0, -1);
    config.params = {};
    config.url = url;
  }
  return config;
};

const requestCatch = (error) => {
  console.error('request err:', error);
  Promise.reject(error);
};

const responseThen = (res: any) => {
  // unset code need to set 200
  const code = res.data['code'] || 200;
  // get error message
  if (code === 200) {
    return Promise.resolve(res.data);
  } else {
    console.error('response err:', res);
    return Promise.reject(res.data);
  }
};

const responseCatch = (error) => {
  console.error('response err:' + error);
  let { message: msg } = error;
  if (msg == 'Network Error') {
    msg = t('requestError.networkError');
  } else if (msg.includes('timeout')) {
    msg = t('requestError.timeout');
  } else if (msg.includes('Request failed with status code')) {
    msg = t('requestError.requestFailed');
  }

  message.error({
    content: msg,
    duration: 5,
  });
  return Promise.reject(error);
};

/**
 * There is no error prompt, you need to implement related error handling yourself
 */
const pureRequest = axios.create(initConfig);

pureRequest.interceptors.request.use(requestThen, requestCatch);
pureRequest.interceptors.response.use(responseThen, responseCatch);

/**
 * when request is error,this func will show error message modal
 */
const getService = () => {
  const service = axios.create(initConfig);
  service.interceptors.request.use(requestThen, requestCatch);
  service.interceptors.response.use((res: any) => {
    return responseThen(res).catch((data) => {
      const code = data['code'] || 200;
      const msg = errorCodeType(code) || data['msg'] || errorCodeType('default');
      message.error(msg);
      return Promise.reject(data);
    });
  }, responseCatch);

  return service;
};

/**
 * If there is an error on the server side, a prompt will pop up automatically
 */
const Request = getService();

export { pureRequest, Request };
