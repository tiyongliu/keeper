import {apiCall} from '/@/second/utility/api'
import { dataBaseStore } from "/@/store/modules/dataBase"
let openedConnectionsHandle: null | number = null

const doServerPing = value => {
  void apiCall('bridge.ServerConnections.Ping', value)
}

export function subscribeConnectionPingers() {
  const dataBase = dataBaseStore()

  dataBase.$subscribe(({events}, state) => {
    if (events.hasOwnProperty('key') && events?.key === 'openedConnections') {
      doServerPing(state.openedConnections)
      if (openedConnectionsHandle) window.clearInterval(openedConnectionsHandle)
      openedConnectionsHandle = window.setInterval(() => doServerPing(state.openedConnections), 30 * 1000)
    }
  })
}
