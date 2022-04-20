import { defineStore } from "pinia";
import { store } from "/@/store";

import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'
import {IConnectionAppObjectData, IPinnedDatabasesItem} from '/@/second/types/IStore.d'
import {ExtensionsDirectory} from '/@/second/types/extensions.d'
interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | {
    name: string
    connection: IConnectionAppObjectData
  },
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
    getPinnedDatabases(): IPinnedDatabasesItem[] {
      return this.pinnedDatabases
    }
  },
  actions: {
    subscribeOpenedConnections(value: string[]) {
      this.openedConnections = value
    },
    subscribeCurrentDatabase(value) {
      this.currentDatabase = value
    },
    subscribeExtensions(value: ExtensionsDirectory) {
      this.extensions = value
    },
    subscribePinnedDatabases(value: IPinnedDatabasesItem[]) {
      this.pinnedDatabases = value
      setWithStorageVariableCache('pinnedDatabases', value)
    }
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}
