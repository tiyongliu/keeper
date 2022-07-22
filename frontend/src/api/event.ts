import {EventsOn} from '/@/wailsjs/runtime/runtime'
import {metadataLoadersStore} from "/@/store/modules/metadataLoaders"
const metadataLoaders = metadataLoadersStore()

export function connectionListChangedEvent() {
  EventsOn('connection-list-changed', async (data) => {
    // const {conid} = data
    // try {
    //   await metadataLoaders.onConnectionGet({conid})
    //   await metadataLoaders.onConnectionList()
    // } catch (e) {
    //   console.log(e)
    // }
  })
}


export function serverStatusChangedEvent() {
  EventsOn('server-status-changed',  _ => {
    void metadataLoaders.onServerStatus()
  })
}


export function databaseListChangedEvent() {
  EventsOn('database-list-changed', conid => {
    void metadataLoaders.onCacheDatabaseList(conid)
  })
}
export function databaseStructureChangedEvent() {
  EventsOn('database-structure-changed', ({conid, database}) => {
    console.log(conid, database, `databaseStructureChangedEvent`)
  })
}
