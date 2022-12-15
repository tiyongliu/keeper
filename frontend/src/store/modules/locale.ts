import type {LocaleSetting, LocaleType} from '/#/config';

import {defineStore} from 'pinia';
import {store} from '/@/store';

import {LOCALE_KEY} from '/@/enums/cacheEnum';
import {createLocalStorage} from '/@/utils/cache';
import {localeSetting} from '/@/settings/localeSetting';
import {getWithStorageVariableCache, setWithStorageVariableCache} from "/@/second/utility/storage";
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {TabDefinition} from "/@/store/modules/bootstrap";
import {isNumber} from "lodash-es";
import {reactive} from "vue";

const ls = createLocalStorage();

const lsLocaleSetting = (ls.get(LOCALE_KEY) || localeSetting) as LocaleSetting;

interface LocaleState {
  localInfo: LocaleSetting;
  selectedWidget: Nullable<string>
  recentDatabases: unknown[]
  pinnedDatabases: IPinnedDatabasesItem[]
  pinnedTables: any[]
  openedTabs: TabDefinition[]
  currentDropDownMenu: null
  leftPanelWidth: number
  visibleTitleBar: number
  dynamicProps: {
    splitterVisible: boolean
  }
}

const LEFTPANELWIDTH = "leftPanelWidth"
export const dynamicProps = reactive({splitterVisible: false})
const _leftPanelWidth = getWithStorageVariableCache(300, LEFTPANELWIDTH)

export const useLocaleStore = defineStore({
  id: 'app-locale',
  state: (): LocaleState => ({
    localInfo: lsLocaleSetting,
    openedTabs: getWithStorageVariableCache<TabDefinition[]>([], 'openedTabs'),
    selectedWidget: getWithStorageVariableCache('database', 'selectedWidget'),
    recentDatabases: getWithStorageVariableCache([], 'recentDatabases'),
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'),
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
    currentDropDownMenu: null,
    visibleTitleBar: 0,
    leftPanelWidth: parseFloat(_leftPanelWidth).toString() !== 'NaN' ?
      parseFloat(_leftPanelWidth) : 300,
    dynamicProps: {
      splitterVisible: false
    }
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
    setPinnedDatabases(value: IPinnedDatabasesItem[]) {
      this.pinnedDatabases = value
      setWithStorageVariableCache('pinnedDatabases', value)
    },
    setPinnedTables(value: any[]) {
      this.pinnedTables = value
      setWithStorageVariableCache('pinnedTables', value)
    },
    setLeftPanelWidth(value) {
      this.leftPanelWidth += value
      document.documentElement.style.setProperty("--dim-left-panel-width", `${this.leftPanelWidth}px`);
      if (isNumber(this.leftPanelWidth)) {
        setWithStorageVariableCache(LEFTPANELWIDTH, String(this.leftPanelWidth));
      }
    },
    setCssVariable(value, transform, cssVariable) {
      document.documentElement.style.setProperty(cssVariable, transform(value));
    },
    setDynamicProps(value: { splitterVisible: boolean }) {
      this.dynamicProps = value
    },
    setCurrentDropDownMenu() {

    },
    setOpenedTabs(value) {
      if (typeof value == 'function') {
        this.openedTabs = value()
      } else {
        this.openedTabs = value
      }
      localStorage.setItem('openedTabs', JSON.stringify(this.openedTabs))
    },
    setRecentDatabases() {

    }
  },
});

// Need to be used outside the setup
export function useLocaleStoreWithOut() {
  return useLocaleStore(store);
}
