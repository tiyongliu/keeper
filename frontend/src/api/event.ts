import {storeToRefs} from 'pinia'
import {EventsOn, EventsEmit}from '/@/wailsjs/runtime/runtime'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {findEngineDriver} from '/@/second/keeper-tools'
import {useClusterApiStore} from '/@/store/modules/clusterApi'
import {dumpSqlSelect} from '/@/second/keeper-sqltree'
import {Select} from '/@/second/keeper-sqltree/types'
import Mongo from '/@/second/plugins/keeper-plugin-mongo'
import Mysql from '/@/second/plugins/keeper-plugin-mysql'

export default function dispatchRuntimeEvent() {
  const bootstrap = useBootstrapStore()
  const {extensions} = storeToRefs(bootstrap)
  const clusterApi = useClusterApiStore()
  const {connection} = storeToRefs(clusterApi)

  //获取数据库基本配置信息
  EventsOn('pullEventPluginsScript', (adapter: 'mongo' | 'mysql') => {
    switch (adapter) {
      case 'mysql':
        EventsEmit('loadPlugins', Mysql)
        break
      case 'mongo':
        EventsEmit('loadPlugins', Mongo)
        break
      default:
        EventsEmit('loadPlugins')
    }
  })

  //构建分页查询条件
  EventsOn('handleSqlSelect', (select: Select) => {
    const driver = extensions.value ? findEngineDriver(connection.value, extensions.value) : null
    if (driver) {
      const dmp = driver?.createDumper()
      dumpSqlSelect(dmp, select)
      EventsEmit('handleSqlSelectReturn', dmp?.s)
    }
  })
}


