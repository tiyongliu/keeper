import {defineStore} from 'pinia'
import {
  getConnectionInfo,
  getConnectionList,
  useServerStatus,
  getDatabaseList,
} from '/@/api/metadataLoaders'
import {TablesNameSort} from '/@/second/typings/mysql'
export enum metadataLoadersKey {
  connections = 'connections',
}

export const metadataLoadersStore = defineStore({
  id: 'app-metadataLoaders',
  state: () => ({
    connections: [],
    databases: []
  }),
  getters: {
    connectionsWithStatus(): unknown[] {
      return this.connections || []
    },
    getDatabaseList(): TablesNameSort[] {
      return this.databases
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
      const serverStatus = await useServerStatus()
      if (this.connections && serverStatus) {
        // @ts-ignore
        this.connections = this.connections.map(conn => ({...conn, status: serverStatus[conn._id]}))
      }
    },
    subscribeDatabaseList(value) {
      this.databases = value
    },
    async onCacheDatabaseList(conid) {
      const data = await getDatabaseList(conid)
      console.log(`onCacheDatabaseList line 54`, data)
      this.databases = data
    }
  }
})
