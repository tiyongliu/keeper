import { defineStore } from "pinia";
import { store } from "/@/store";

import pluginMongoDrivers from '/@/second/plugins/keeper-plugin-mongo'
import pluginMysqlDrivers from '/@/second/plugins/keeper-plugin-mysql'
import pluginRedisDrivers from '/@/second/plugins/keeper-plugin-redis'

import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {ExtensionsDirectory} from '/@/second/typings/types/extensions.d'
interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | IPinnedDatabasesItem,
  extensions: null | ExtensionsDirectory
  pinnedDatabases: IPinnedDatabasesItem[],
  pinnedTables: []
}

export const dataBaseStore = defineStore({
  id: "app-dataBase",
  state: (): IVariableBasic => ({
    currentDatabase: null,
    openedConnections: [],
    extensions: {
      drivers: [].concat(pluginMongoDrivers.drivers).concat(pluginMysqlDrivers.drivers).concat(pluginRedisDrivers.drivers)
    },
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'),
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
  }),
  getters: {
    getCurrentDatabase(): IPinnedDatabasesItem | null {
      return this.currentDatabase
    },
    getPinnedDatabases(): IPinnedDatabasesItem[] {
      return this.pinnedDatabases
    },
    getPinnedTables(): [] {
      return this.pinnedTables
    },
    getPinnedExtensions() {
      return this.extensions
    }
  },
  actions: {
    subscribeOpenedConnections(value: string[]) {
      this.openedConnections = value
    },
    subscribeCurrentDatabase(value: IPinnedDatabasesItem) {
      this.currentDatabase = value
    },
    subscribeExtensions(value: ExtensionsDirectory) {
      this.extensions = value
    },
    subscribePinnedDatabases(value: IPinnedDatabasesItem[]) {
      this.pinnedDatabases = value
      setWithStorageVariableCache('pinnedDatabases', this.pinnedDatabases)
    },
    subscribePinnedTables(value: any) {
      this.pinnedTables = value
    }
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}
