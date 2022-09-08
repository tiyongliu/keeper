import {apiCall} from '/@/second/utility/api'

export async function handleScriptApi(params: {packageName: string}) {
  return await apiCall('bridge.Plugins.Script', params)
}
