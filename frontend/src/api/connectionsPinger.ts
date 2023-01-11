import {watch} from 'vue'
import {get} from 'lodash-es'
import {storeToRefs} from 'pinia'
import {apiCall} from '/@/second/utility/api'
import {useBootstrapStore} from "/@/store/modules/bootstrap"

let openedConnectionsHandle: null | number = null
let currentDatabaseHandle: null | number = null

const doServerPing = value => {
  void apiCall('bridge.ServerConnections.Ping', {connections: value})
}

const doDatabasePing = value => {
  const database = get(value, 'name')
  const conid = get(value, 'connection._id')
  if (conid && database) {
    void apiCall('bridge.DatabaseConnections.Ping', { conid, database })
  }
}

export function subscribeConnectionPingers() {
  const bootstrap = useBootstrapStore()
  const {openedConnections, currentDatabase} = storeToRefs(bootstrap)
  watch(() => openedConnections.value, () => {
    doServerPing(openedConnections.value)
    if (openedConnectionsHandle) window.clearInterval(openedConnectionsHandle)
    openedConnectionsHandle = window.setInterval(() => doServerPing(openedConnections.value), 30 * 1000)
  }, {immediate: true})

  watch(() => currentDatabase.value, () => {
    doDatabasePing(currentDatabase.value)
    if (currentDatabaseHandle) window.clearInterval(currentDatabaseHandle)
    currentDatabaseHandle = window.setInterval(() => doDatabasePing(currentDatabase.value), 30 * 1000)
  }, {immediate: true})
}
