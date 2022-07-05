import {apiCall} from "/@/second/utility/api";

export async function handleRefreshApi(params) {
  return await apiCall('bridge.ServerConnections.Refresh', params)
}
