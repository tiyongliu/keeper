import {useGlobSetting} from '/@/hooks/setting';
import {defHttp} from '/@/utils/http/axios';

const {environment} = useGlobSetting()

export async function apiCall<T>(relativePath: string, params?: T): Promise<T | void> {
  //读取环境变量
  console.log('>>> API CALL', relativePath, params)

  if (environment === 'web') {
    return await defHttp.post({url: relativePath, params})
  } else {

    // window['go']['proc']

    // const resp = await window['go']['backend']['MMMM']['GetVersion'](params)
    // window['go']['proc']['ConnectProcess']['Test']()
  }
}
