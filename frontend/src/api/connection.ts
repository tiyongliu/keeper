import {apiCall} from '/@/second/utility/api'

/**
* @description 测试连接池
* @param {Object} params 连接信息
*/

var a = {
  relativePath: '/connections/test',
  params: function () {}
}


export async function handleDriverTestApi<T>(params) {
  const resp = await apiCall<T>('/connections/test', params)
}
