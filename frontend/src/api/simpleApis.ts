import {apiCall} from '/@/second/utility/api'

/**
 * @description connection
 */
export async function connectionTestApi(params) {
  return await apiCall('bridge.Connections.Test', params)
}

export async function connectionSaveApi(params) {
  return await apiCall('bridge.Connections.Save', params)
}

export async function connectionListApi() {
  return await apiCall('bridge.Connections.List')
}

export async function connectionDeleteApi(params) {
  return await apiCall('bridge.Connections.Delete', params)
}

/**
 * @description databaseConnections
 */
export async function databaseConnectionsRefreshApi(params: { conid: string, database: string, keepOpen?: boolean }) {
  return await apiCall('bridge.DatabaseConnections.Refresh', params)
}

export async function databaseConnectionsSqlSelectApi<T>(params: { conid: string, database: string, select: unknown }): Promise<T> {
  return await apiCall('bridge.DatabaseConnections.SqlSelect', params) as T
}

export async function databaseConnectionsCollectionDataApi<T>(params: { conid: string, database: string, options: unknown }): Promise<T> {
  return await apiCall('bridge.DatabaseConnections.CollectionData', params) as T
}

/**
 * @description plugins
 */
export async function pluginsScriptApi(params: { packageName: string }) {
  return await apiCall('bridge.Plugins.Script', params)
}

/**
 * @description serverConnections
 */
export async function serverConnectionsRefreshApi(params) {
  return await apiCall('bridge.ServerConnections.Refresh', params)
}
