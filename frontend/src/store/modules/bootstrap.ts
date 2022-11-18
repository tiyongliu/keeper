import {Component, PropType} from 'vue'
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
  connections: []
  databases: []
  selectedCellsCallback: Nullable<() => any>
}

export interface TabDefinition {
  title: string;
  closedTime?: number;
  icon: string;
  props: any;
  selected: boolean;
  busy: boolean;
  tabid: string;
  tabComponent: PropType<string | Component>;
  tabOrder?: number;
}

export interface ICurrentDropDownMenu {
  left: number
  top: number
  items: any[]
  targetElement?: HTMLElement
}

let visibleCommandPaletteValue = null
export const useBootstrapStore = defineStore({
  id: "app-bootstrap",
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
    },
    connections: [],
    databases: [],
    selectedCellsCallback: null
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
    },
    getCommandsCustomized(): any[] {
      return mapValues([this.commands, this.commandsSettings], (v, k) => ({
        // @ts-ignore
        ...v,
        ...this.commandsSettings[k]
      }))
    }
  },
  actions: {
    subscribeOpenedConnections(value: string[]) {
      this.openedConnections = value
    },
    subscribeCurrentDatabase(value: null | IPinnedDatabasesItem) {
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
      console.log(visibleCommandPaletteValue)
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
    },
    subscribeConnections(payload): void {
      this.connections = payload
    },
    subscribeSelectedCellsCallback(value: () => any) {
      this.selectedCellsCallback = value
    }
  }
});

export function useBootstrapStoreWithOut() {
  return useBootstrapStore(store);
}

const derived = (commands, commandsSettings): object => {
  return mapValues(commands, (v, k) => ({
    ...v,
    ...commandsSettings[k]
  }))
}
