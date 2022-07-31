import {useGlobSetting} from '/@/hooks/setting';

const {environment} = useGlobSetting()

const apiLogging = false
export async function apiCall<T>(url: string, params?: T): Promise<T | void> {
  //读取环境变量
  if (apiLogging) {
    console.log('>>> API CALL', url, params)
  }

  if (environment === 'web') {
    //TODO 暂时不支持http方式访问
    // return await defHttp.post({url: relativePath, params})
  } else {
    try {
      let self: Function = window['go'];
      url.split(/[.\/]/).filter(item => item).forEach(key => self = self[key])
      if (!params || Object.keys(params).length === 0) {
        const resp = await self()
        return processApiResponse(url, params, resp)
      }
      const resp = await self(params)
      return processApiResponse(url, params, resp)
    } catch (e) {
      return Promise.reject(e)
    }
  }
}

function processApiResponse(url, params, resp) {
  if (!apiLogging) {
    console.log('<<< API RESPONSE', url, params, resp)
  }

  if (resp.status === 1) {
    return {
      errorMessage: resp.message
    }
    // return resp.result.message
    // throw resp.message
  }
  return resp.result
}
