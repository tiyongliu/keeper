import {useGlobSetting} from '/@/hooks/setting';

const {environment} = useGlobSetting()

let apiLogging = true

export async function apiCall<T>(relativePath: string, params?: T): Promise<T | void> {
  //读取环境变量
  if (apiLogging) {
    console.log('>>> API CALL', relativePath, params)
  }

  if (environment === 'web') {
    //todo 暂时不支持http方式访问
    // return await defHttp.post({url: relativePath, params})
  } else {
    try {
      let self: Function = window['go'];
      relativePath.split(/[.|\/]/).filter(item => item).forEach(key => self = self[key])
      if (!params) {
        const resp = await self()
        return processApiResponse(relativePath, params, resp)
      }
      const resp = await self(params)
      return processApiResponse(relativePath, params, resp)
    } catch (e) {
      return Promise.reject(e)
    }
  }
}

function processApiResponse(relativePath, params, resp) {
  if (apiLogging) {
    console.log('<<< API RESPONSE', relativePath, params, resp)
  }

  if (resp.code === 1) {
    // return resp.result.message
    return
  }
  return resp.result
}
