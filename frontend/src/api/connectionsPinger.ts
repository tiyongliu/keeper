import {get} from 'lodash-es'
import {apiCall} from '/@/second/utility/api'
import {useBootstrapStore} from "/@/store/modules/bootstrap"

let openedConnectionsHandle: null | number = null
let currentDatabaseHandle: null | number = null

const doServerPing = value => {
  void apiCall('bridge.ServerConnections.Ping', value)
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
  bootstrap.$subscribe((mutation, state) => {
    const {events} = mutation as any
    if (events && events.hasOwnProperty('key') && events.key === 'openedConnections') {
      doServerPing(state.openedConnections)
      if (openedConnectionsHandle) window.clearInterval(openedConnectionsHandle)
      openedConnectionsHandle = window.setInterval(() => doServerPing(state.openedConnections), 30 * 1000)
    }

    if (events && events.hasOwnProperty('key') && events.key === 'currentDatabase') {
      doDatabasePing(state.currentDatabase)
      if (currentDatabaseHandle) window.clearInterval(currentDatabaseHandle)
      currentDatabaseHandle = window.setInterval(() => doDatabasePing(state.currentDatabase), 30 * 1000)
    }
  })
}
