import {apiCall} from '/@/second/utility/api'

export async function handleRefreshApi(params: {conid: string, database: string}) {
  return await apiCall('bridge.DatabaseConnections.Refresh', params)
}
