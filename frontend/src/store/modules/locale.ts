import type {LocaleSetting, LocaleType} from '/#/config';

import {defineStore} from 'pinia';
import {store} from '/@/store';

import {LOCALE_KEY} from '/@/enums/cacheEnum';
import {createLocalStorage} from '/@/utils/cache';
import {localeSetting} from '/@/settings/localeSetting';
import {getWithStorageVariableCache, setWithStorageVariableCache} from "/@/second/utility/storage";
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {TabDefinition} from "/@/store/modules/dataBase";

const ls = createLocalStorage();

const lsLocaleSetting = (ls.get(LOCALE_KEY) || localeSetting) as LocaleSetting;

interface LocaleState {
  localInfo: LocaleSetting;
  selectedWidget: null | string
  pinnedDatabases: IPinnedDatabasesItem[]
  pinnedTables: any[]
  openedTabs: any[]
}

export const useLocaleStore = defineStore({
  id: 'app-locale',
  state: (): LocaleState => ({
    localInfo: lsLocaleSetting,
    // selectedWidget: null
    openedTabs: getWithStorageVariableCache<TabDefinition[]>([], 'openedTabs'),
    selectedWidget: getWithStorageVariableCache('database', 'selectedWidget'),
    pinnedDatabases: getWithStorageVariableCache([], 'pinnedDatabases'),
    pinnedTables: getWithStorageVariableCache([], 'pinnedTables'),
  }),
  getters: {
    getShowPicker(): boolean {
      return !!this.localInfo?.showPicker;
    },
    getLocale(): LocaleType {
      return this.localInfo?.locale ?? 'zh_CN';
    },
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
    setSelectedWidget(name: string | null) {
      this.selectedWidget = name
      setWithStorageVariableCache('selectedWidget', name)
    },
    subscribePinnedDatabases(value: IPinnedDatabasesItem[]) {
      this.pinnedDatabases = value
      setWithStorageVariableCache('pinnedDatabases', this.pinnedDatabases)
    },
    subscribePinnedTables(value: any[]) {
      this.pinnedTables = value
    },
  },
});

// Need to be used outside the setup
export function useLocaleStoreWithOut() {
  return useLocaleStore(store);
}
