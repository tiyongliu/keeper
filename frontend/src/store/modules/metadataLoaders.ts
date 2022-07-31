import {defineStore} from 'pinia'
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
  }
})
