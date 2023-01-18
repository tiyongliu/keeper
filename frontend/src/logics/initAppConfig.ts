/**
 * Application configuration
 */
import type {ProjectConfig} from '/#/config';
import {PROJ_CFG_KEY} from '/@/enums/cacheEnum';
import projectSetting from '/@/settings/projectSetting';
import {updateDarkTheme} from '/@/logics/theme/dark';
import {useAppStore} from '/@/store/modules/app';
import {getCommonStoragePrefix, getStorageShortName} from '/@/utils/env';
import {Persistent} from '/@/utils/cache/persistent';
import {deepMerge} from '/@/utils';
import {ThemeEnum} from '/@/enums/appEnum'

// Initial project configuration
export function initAppConfigStore() {
  const appStore = useAppStore();
  appStore.setDarkMode(ThemeEnum.LIGHT) // todo 主题色切换核心
  let projCfg: ProjectConfig = Persistent.getLocal(PROJ_CFG_KEY) as ProjectConfig;
  projCfg = deepMerge(projectSetting, projCfg || {});
  const darkMode = appStore.getDarkMode;
  appStore.setProjectConfig(projCfg);

  // init dark mode
  updateDarkTheme(darkMode);
  setTimeout(() => {
    clearObsoleteStorage();
  }, 16);
}

/**
 * As the version continues to iterate, there will be more and more cache keys stored in localStorage.
 * This method is used to delete useless keys
 */
export function clearObsoleteStorage() {
  const commonPrefix = getCommonStoragePrefix();
  const shortPrefix = getStorageShortName();

  [localStorage, sessionStorage].forEach((item: Storage) => {
    Object.keys(item).forEach((key) => {
      if (key && key.startsWith(commonPrefix) && !key.startsWith(shortPrefix)) {
        item.removeItem(key);
      }
    });
  });
}
