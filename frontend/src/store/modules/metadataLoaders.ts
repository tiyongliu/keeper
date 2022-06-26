import {defineStore} from 'pinia'
import {
  getConnectionInfo,
  getConnectionList,
  useServerStatus
} from '/@/second/utility/metadataLoaders'

export enum metadataLoadersKey {
  connections = 'connections',
}

export const metadataLoadersStore = defineStore({
  id: 'app-metadataLoaders',
  state: () => ({
    connections: []
  }),
  getters: {
    connectionsWithStatus(): unknown[] {
      return this.connections || []
    }
  },
  actions: {
    setState<T>(type: string, payload: T): void {
      this[type] = payload
    },
    async onConnectionGet(args) {
      await getConnectionInfo(args)
    },
    async onConnectionList() {
      this.connections = await getConnectionList()
    },
    setConnectionList(value) {
      this.connections = value
    },
    async onServerStatus() {
      const data = await useServerStatus()
      console.log(`fnksdnfksnksg`, data)
    },
  }
})
