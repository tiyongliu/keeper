import {apiCall} from '/@/second/utility/api'
import { dataBaseStore } from "/@/store/modules/dataBase"

const doServerPing = value => {
  console.log(`bridge.ServerConnections.Ping`, value)
  apiCall('bridge.ServerConnections.Ping', { connections: value })
    .then(res => console.log(res, `res`)).catch(err => console.log(err, `err`))
  // apiCall('server-connections/ping', { connections: value });
}

// const doDatabasePing = value => {
//   const database = get(value, 'name')
//   const conid = get(value, 'connection._id')
//   if (conid && database) {
//     apiCall('bridge.DatabaseConnections.Ping', {
//       conid,
//       database
//     }).then(res => console.log(res, `res`)).catch(err => console.log(err, `err`))
//   }
// }

let openedConnectionsHandle: null | number = null

export function subscribeConnectionPingers() {
  const dataBase = dataBaseStore()

  dataBase.$subscribe((_, state) => {
    doServerPing(state.openedConnections)
    if (openedConnectionsHandle) window.clearInterval(openedConnectionsHandle)
    openedConnectionsHandle = window.setInterval(() => doServerPing(state.openedConnections), 30 * 1000)
    // console.log(mutation)
  })
}
