import {defineStore} from 'pinia'

export const metadataLoadersStore = defineStore({
  id: 'app-metadataLoaders',
  state: () => ({
    connections: [],
    databases: []
  }),
  actions: {
    setConnections(payload): void {
      this.connections = payload
    },
  }
})
