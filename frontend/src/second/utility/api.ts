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
    let self: Function = window['go']['proc'];
    relativePath.split(/[.|\/]/).filter(item => item).forEach(key => self = self[key])
    const resp = await self(params)
    return processApiResponse(relativePath, params, resp)
  }
}

function processApiResponse(relativePath, params, resp) {
  if (apiLogging) {
    console.log('<<< API RESPONSE', relativePath, params, resp)
  }

  if (resp.code === 0) {
    return resp.result.message
  }
  return resp.result
}
