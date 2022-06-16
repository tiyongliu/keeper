import {EventsOn} from '/@/wailsjs/runtime/runtime'
import {metadataLoadersStore} from "/@/store/modules/metadataLoaders"
const metadataLoaders = metadataLoadersStore()

export function connectionListChangedEvent() {
  EventsOn('connection-list-changed', async (data) => {
    const {conid} = data
    await metadataLoaders.onConnectionGet({conid})
    await metadataLoaders.onConnectionList()
  })
}
