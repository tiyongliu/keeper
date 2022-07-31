import {defineStore} from "pinia"
import {store} from "/@/store"
import {mapValues} from 'lodash-es'
import invalidateCommands from '/@/second/commands/invalidateCommands'
import pluginMongoDrivers from '/@/second/plugins/keeper-plugin-mongo'
import pluginMysqlDrivers from '/@/second/plugins/keeper-plugin-mysql'
import pluginRedisDrivers from '/@/second/plugins/keeper-plugin-redis'
import {apiCall} from '/@/second/utility/api'
import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {ExtensionsDirectory} from '/@/second/typings/types/extensions.d'

interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | IPinnedDatabasesItem,
  extensions: null | ExtensionsDirectory
  pinnedDatabases: IPinnedDatabasesItem[],
  pinnedTables: [],
  currentDropDownMenu: null | ICurrentDropDownMenu
  commands: object
  commandsSettings: object
  visibleCommandPalette: null | unknown
  commandsCustomized: object
}

export interface TabDefinition {
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

export interface ICurrentDropDownMenu {
  left: number
  top: number
  items: any[]
  targetElement: HTMLElement
}

let visibleCommandPaletteValue = null
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
    openedTabs: getWithStorageVariableCache<TabDefinition[]>([], 'openedTabs'),
    currentDropDownMenu: null,
    commands: {},
    commandsSettings: {},
    visibleCommandPalette: null,
    commandsCustomized: {},
  }),
  getters: {
    getOpenedConnections(): string[] {
      return this.openedConnections
    },
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
    },
    getPinnedExtensions() {
      return this.extensions
    }
  },
  actions: {
    setState<T>(type: string, payload: T): void {
      this[type] = payload
    },
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
    },
    subscribeCurrentDropDownMenu(value: null | ICurrentDropDownMenu) {
      this.currentDropDownMenu = value
    },
    subscribeVisibleCommandPalette(value) {
      visibleCommandPaletteValue = value
      void invalidateCommands()
    },
    setVisibleCommandPalette(value: null | unknown) {
      this.visibleCommandPalette = value
    },
    subscribeCommands(value: object) {
      this.commands = value

      this.commandsCustomized = derived(this.commands, this.commandsSettings)
    },
    subscribeCommandsSettings(value: object) {
      this.commandsSettings = value
      this.commandsCustomized = derived(this.commands, this.commandsSettings)
    },
  }
});

export function useDataBaseStoreWithOut() {
  return dataBaseStore(store);
}

const derived = (commands, commandsSettings): object => {
  return mapValues(commands, (v, k) => ({
    ...v,
    ...commandsSettings[k]
  }))
}
