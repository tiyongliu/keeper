import {reactive} from 'vue'
import {defineStore} from 'pinia'
import {isNumber} from 'lodash-es'
import type {LocaleSetting, LocaleType} from '/#/config'
import {store} from '/@/store'
import invalidateCommands from '/@/second/commands/invalidateCommands'
import {LOCALE_KEY} from '/@/enums/cacheEnum'
import {getWithStorageVariableCache, setWithStorageVariableCache} from '/@/second/utility/storage'
import {TabDefinition} from '/@/store/modules/bootstrap'
import {localeSetting} from '/@/settings/localeSetting'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {createLocalStorage} from '/@/utils/cache'


const ls = createLocalStorage();

const lsLocaleSetting = (ls.get(LOCALE_KEY) || localeSetting) as LocaleSetting;

interface LocaleState {
  localInfo: LocaleSetting;
  selectedWidget: Nullable<string>
  recentDatabases: unknown[]
  pinnedDatabases: IPinnedDatabasesItem[]
  pinnedTables: any[]
  openedTabs: TabDefinition[]
  leftPanelWidth: number
  visibleTitleBar: number
  dynamicProps: {
    splitterVisible: boolean
  },
  copyRowsFormat: string
}

const LEFTPANELWIDTH = "leftPanelWidth"
export const dynamicProps = reactive({splitterVisible: false})
const _leftPanelWidth = getWithStorageVariableCache(300, LEFTPANELWIDTH)
let openedTabsValue: Nullable<TabDefinition[]> = null
export const getOpenedTabs = () => openedTabsValue
export const useLocaleStore = defineStore({
  id: 'app-locale',
  state: (): LocaleState => ({
    localInfo: lsLocaleSetting,
    openedTabs: getWithStorageVariableCache<TabDefinition[]>([], 'openedTabs'),
    selectedWidget: getWithStorageVariableCache('database', 'selectedWidget'),
    recentDatabases: getWithStorageVariableCache([], 'recentDatabases'),
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'),
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
    visibleTitleBar: 0,
    leftPanelWidth: parseFloat(_leftPanelWidth).toString() !== 'NaN' ?
      parseFloat(_leftPanelWidth) : 300,
    dynamicProps: {
      splitterVisible: false
    },
    copyRowsFormat: getWithStorageVariableCache('textWithoutHeaders', 'copyRowsFormat')
  }),
  getters: {
    getShowPicker(): boolean {
      return !!this.localInfo?.showPicker;
    },
    getLocale(): LocaleType {
      return this.localInfo?.locale ?? 'zh_CN';
    },
    getDynamicProps(): { splitterVisible: boolean } {
      return this.dynamicProps
    },
    getOpenedTabs(): TabDefinition[] {
      return this.openedTabs
    },
    activeTabId(): string | undefined {
      return this.openedTabs.find(x => x.selected)?.tabid
    },
    activeTab(): any {
      return this.openedTabs.find(x => x.selected)
    }
  },
  actions: {
    /**
     * Set up multilingual information and cache
     * @param info multilingual info
     */
    setLocaleInfo(info: Partial<LocaleSetting>) {
      this.localInfo = {...this.localInfo, ...info};
      ls.set(LOCALE_KEY, this.localInfo);
    },
    /**
     * Initialize multilingual information and load the existing configuration from the local cache
     */
    initLocale() {
      this.setLocaleInfo({
        ...localeSetting,
        ...this.localInfo,
      });
    },
    setSelectedWidget(name: Nullable<string>) {
      this.selectedWidget = name
      setWithStorageVariableCache('selectedWidget', name)
    },
    updatePinnedDatabases(updater: (list: IPinnedDatabasesItem[]) => IPinnedDatabasesItem[]) {
      this.pinnedDatabases = updater(this.pinnedDatabases)
      setWithStorageVariableCache('pinnedDatabases', this.pinnedDatabases)
    },
    updatePinnedTables(updater: (list: any[]) => any[]) {
      this.pinnedTables = updater(this.pinnedTables)
      setWithStorageVariableCache('pinnedTables', this.pinnedTables)
    },
    updateLeftPanelWidth(updater) {
      this.leftPanelWidth = updater(this.leftPanelWidth)
      document.documentElement.style.setProperty('--dim-left-panel-width', `${this.leftPanelWidth}px`);
      if (isNumber(this.leftPanelWidth)) {
        setWithStorageVariableCache(LEFTPANELWIDTH, String(this.leftPanelWidth));
      }
    },
    setCopyRowsFormat(value: string) {
      this.copyRowsFormat = value
      setWithStorageVariableCache('copyRowsFormat', value)
    },
    setCssVariable(value, transform, cssVariable) {
      document.documentElement.style.setProperty(cssVariable, transform(value));
    },
    setDynamicProps(value: { splitterVisible: boolean }) {
      this.dynamicProps = value
    },
    updateOpenedTabs(updater) {
      this.openedTabs = updater(this.openedTabs)
      setWithStorageVariableCache('openedTabs', this.openedTabs)
      openedTabsValue = this.openedTabs
      void invalidateCommands()
    },
    updateRecentDatabases(updater) {
      this.recentDatabases = updater(this.recentDatabases)
      setWithStorageVariableCache('recentDatabases', this.recentDatabases)
    }
  },
});

// Need to be used outside the setup
export function useLocaleStoreWithOut() {
  return useLocaleStore(store);
}
