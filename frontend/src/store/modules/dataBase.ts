import {defineStore} from "pinia"
import {store} from "/@/store"
import {mapValues} from 'lodash-es'
import invalidateCommands from '/@/second/commands/invalidateCommands'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {ExtensionsDirectory} from '/@/second/typings/types/extensions.d'

interface IVariableBasic {
  openedConnections: string[]
  currentDatabase: null | IPinnedDatabasesItem,
  extensions: ExtensionsDirectory | null
  currentDropDownMenu: null | ICurrentDropDownMenu
  commands: object
  commandsSettings: object
  visibleCommandPalette: null | unknown
  commandsCustomized: object
  loadingPluginStore: { loaded: boolean, loadingPackageName: string | null }
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
  targetElement?: HTMLElement
}

let visibleCommandPaletteValue = null
export const dataBaseStore = defineStore({
  id: "app-dataBase",
  state: (): IVariableBasic => ({
    currentDatabase: null,
    openedConnections: [],
    extensions: null,
    currentDropDownMenu: null,
    commands: {},
    commandsSettings: {},
    visibleCommandPalette: null,
    commandsCustomized: {},
    loadingPluginStore: {
      loaded: true,
      loadingPackageName: null
    }
  }),
  getters: {
    getOpenedConnections(): string[] {
      return this.openedConnections
    },
    getCurrentDatabase(): IPinnedDatabasesItem | null {
      return this.currentDatabase
    },
    getPinnedExtensions(): ExtensionsDirectory | null {
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
    subscribeLoadingPluginStore(value: { loaded: boolean, loadingPackageName: string | null }) {
      this.loadingPluginStore = value
    }
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
