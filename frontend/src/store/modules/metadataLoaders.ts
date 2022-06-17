import { defineStore, mapState } from 'pinia'
import {getConnectionList, getConnectionInfo} from '/@/second/utility/metadataLoaders'

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
    async onConnectionGet(args) {
      await getConnectionInfo(args)
    },
    async onConnectionList() {
      this.connections = await getConnectionList()
    },
    async setConnectionList(value) {
      this.connections = value
    },
    async onServerStatus() {}
  }
})
// connectionsWithStatus
