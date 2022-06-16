import {apiCall} from '/@/second/utility/api'

/**
* @description 测试连接池
* @param {Object} params 连接信息
*/
export async function handleDriverTestApi(params) {
  return await apiCall('bridge.Connections.Test', params)
}

export async function handleDriverSaveApi(params) {
  return await apiCall('bridge.Connections.Save', params)
}

export async function loadDatabasesApi() {
  return await apiCall('bridge.Connections.List')
}
