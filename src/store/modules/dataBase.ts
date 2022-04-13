import { defineStore } from "pinia";
import { store } from "/@/store";

import {IConnectionAppObjectData} from '/@/second/types/IStore.d'
import {ExtensionsDirectory} from '/@/second/types/extensions.d'
interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | {
    name: string
    connection: IConnectionAppObjectData
  },
  extensions: null | ExtensionsDirectory
}

export const dataBaseStore = defineStore({
  id: "app-dataBase",
  state: (): IVariableBasic => ({
    currentDatabase: null,
    openedConnections: [],
    extensions: null
  }),
  getters: {},
  actions: {
    subscribeOpenedConnections(value: string[]) {
      this.openedConnections = value
    },
    subscribeCurrentDatabase(value) {
      this.currentDatabase = value
    },
    subscribeExtensions(value: ExtensionsDirectory) {
      this.extensions = value
    }
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}
