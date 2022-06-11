import { defineStore } from 'pinia'
import {useCore} from '/@/second/utility/metadataLoaders'

const connectionListLoader = () => ({
  url: 'bridge.Connections.List',
  params: null,
  reloadTrigger: `connection-list-changed`
})

const serverStatusLoader = () => ({
  url: 'bridge.Connections.ServerStatus',
  params: {},
  reloadTrigger: `server-status-changed`,
})

export const metadataLoadersStore = defineStore({
  id: 'app-metadataLoaders',
  state: () => ({
    connections: []
  }),
  getters: {
    connectionsWithStatus(): never[]{
      return this.connections
    }
  },
  actions: {
    async useConnectionList() {
      this.connections = await useCore(connectionListLoader, {})
      return this.connections
    },
    async useServerStatus() {
      await useCore(serverStatusLoader, {})
    }
  }
})
// connectionsWithStatus
