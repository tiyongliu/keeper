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

// const dta = [{"name":"mysql","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"performance_schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"kb-dms","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"mallplusbak","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}}]

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
