import { defineStore } from "pinia";
import { store } from "/@/store";

import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {ExtensionsDirectory} from '/@/second/typings/types/extensions.d'
interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | IPinnedDatabasesItem,
  extensions: null | ExtensionsDirectory
  pinnedDatabases: IPinnedDatabasesItem[],
  openedTabs: IPinnedDatabasesItem[],
  pinnedTables: []
}

export interface TabDefinition{
  title: string;
  closedTime?: number;
  icon: string;
  props: any;
  selected: boolean;
  busy: boolean;
  tabid: string;
  tabComponent: string;
  tabOrder?: number;
}

export const dataBaseStore = defineStore({
  id: "app-dataBase",
  state: (): IVariableBasic => ({
    currentDatabase: null,
    openedConnections: [],
    extensions: null,
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'), // writableWithStorage
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
    openedTabs: getWithStorageVariableCache<TabDefinition[]>([], 'openedTabs')
  }),
  getters: {
    getCurrentDatabase(): IPinnedDatabasesItem | null {
      return this.currentDatabase
    },
    getPinnedDatabases(): IPinnedDatabasesItem[] {
      return this.pinnedDatabases
    },
    getOpenedTabs(): IPinnedDatabasesItem[] {
      return this.openedTabs
    },
    getPinnedTables(): [] {
      return this.pinnedTables
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
    subscribeOpenedTabs(value: IPinnedDatabasesItem[]) {
      this.openedTabs = value
      setWithStorageVariableCache('openedTabs', this.openedTabs)
    },
    subscribePinnedTables(value: any) {
      this.pinnedTables = value
    }
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}
// export function writableWithStorage<T>(defaultValue: T,storageName){
//   const init  = localStorage.getItem(storageName);
//   // const res =
// }
