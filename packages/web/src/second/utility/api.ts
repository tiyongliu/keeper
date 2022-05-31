import {useGlobSetting} from '/@/hooks/setting';
import {defHttp} from '/@/utils/http/axios';

const {drivers, environment} = useGlobSetting()

export async function apiCall<T>(relativePath: string, params?: T): Promise<T> {
  //读取环境变量
  console.log('>>> API CALL', relativePath, params)

  if (environment === 'web') {
    const resp = await defHttp.post({url: relativePath, params})
    console.log('<<< API RESPONSE', relativePath, params, resp)

    return resp
  }
  return
}