const apiLogging = false
export async function apiCall<T>(url: string, params?: any): Promise<T | void> {
  //读取环境变量
  if (apiLogging) {
    console.log('>>> API CALL', url, params)
  }

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

function processApiResponse(url, params, resp) {
  if (apiLogging) {
    console.log('<<< API RESPONSE', url, params, resp)
  }

  if (resp.status === 1) {
    return {
      ...resp,
      errorMessage: resp.message
    }
    // return resp.result.message
    // throw resp.message
  }
  return resp.result
}
