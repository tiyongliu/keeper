import {storeToRefs} from 'pinia'
import {EventsOn, EventsEmit}from '/@/wailsjs/runtime/runtime'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {findEngineDriver} from '/@/second/keeper-tools'
import {useClusterApiStore} from '/@/store/modules/clusterApi'
import {dumpSqlSelect} from '/@/second/keeper-sqltree'
import {Select} from '/@/second/keeper-sqltree/types'

export default function dispatchRuntimeEvent() {
  const bootstrap = useBootstrapStore()
  const {extensions} = storeToRefs(bootstrap)
  const clusterApi = useClusterApiStore()
  const {connection} = storeToRefs(clusterApi)

  EventsOn('handleSqlSelect', (select: Select) => {
    const driver = extensions.value ? findEngineDriver(connection.value, extensions.value) : null
    if (driver) {
      const dmp = driver?.createDumper()
      dumpSqlSelect(dmp, select)
      EventsEmit('handleSqlSelectReturn', dmp?.s)
    }
  })
}


