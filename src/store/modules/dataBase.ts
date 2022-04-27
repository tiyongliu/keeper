import { defineStore } from "pinia";
import { store } from "/@/store";

import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'
import {IPinnedDatabasesItem} from '/@/second/types/standard.d'
import {ExtensionsDirectory} from '/@/second/types/extensions.d'
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
    extensions: null,
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'),
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
  }),
  getters: {
    getCurrentDatabase(): IPinnedDatabasesItem | null {
      return this.currentDatabase
    },
    getPinnedDatabases(): IPinnedDatabasesItem[] {
      return this.pinnedDatabases
    }
  },
  actions: {
    subscribeOpenedConnections(value: string[]) {
      console.log(value, `value`)
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
    }
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}
