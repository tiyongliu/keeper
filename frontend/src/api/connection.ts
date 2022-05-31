import {apiCall} from '/@/second/utility/api'

/**
* @description 测试连接池
* @param {Object} params 连接信息
*/
export async function handleDriverTestApi(params) {
  return await apiCall('Connections.Test', params)
}
