import {unref} from 'vue'
import {findIndex, findLastIndex, keys, max, pick, sortBy} from 'lodash-es'
import localforage from 'localforage'
import {storeToRefs} from 'pinia'
import stableStringify from 'json-stable-stringify'
import {useLocaleStore} from '/@/store/modules/locale'
import tabs from '/@/second/tabs'
import {setSelectedTabFunc} from '/@/second/utility/common'
import {buildUUID} from '/@/utils/uuid'

const locale = useLocaleStore()
const {openedTabs: oldTabs} = storeToRefs(locale)

function findFreeNumber(numbers: number[]) {
  if (numbers.length == 0) return 1;
  return max(numbers)! + 1;
  // let res = 1;
  // while (numbers.includes(res)) res += 1;
  // return res;
}

export function getTabDbKey(tab) {
  if (tab.props && tab.props.conid && tab.props.database) {
    return `database://${tab.props.database}-${tab.props.conid}`;
  }
  if (tab.props && tab.props.conid) {
    return `server://${tab.props.conid}`;
  }
  if (tab.props && tab.props.archiveFolder) {
    return `archive://${tab.props.archiveFolder}`;
  }
  return null;
}

export function sortTabs(tabs: any[]): any[] {
  return sortBy(tabs, [x => x.tabOrder || 0, x => getTabDbKey(x), 'title', 'tabid']);
}

export function groupTabs(tabs: any[]) {
  const res: any = [];

  for (const tab of sortTabs(tabs)) {
    const lastGroup = res[res.length - 1];
    if (lastGroup && tab.tabDbKey && lastGroup.tabDbKey == tab.tabDbKey) {
      lastGroup.tabs.push(tab);
    } else {
      res.push({
        tabDbKey: tab.tabDbKey,
        tabDbName: tab.tabDbName,
        tabs: [tab],
        grpid: tab.tabid,
      });
    }
  }

  return res;
}

export default async function openNewTab(newTab, initialData: any = undefined, options: unknown = undefined) {
  let existing: unknown = null;

  const {savedFile, savedFolder, savedFilePath} = newTab.props || {}
  if (savedFile || savedFilePath) {
    existing = unref(oldTabs).find(
      x =>
        x.props &&
        x.tabComponent == newTab.tabComponent &&
        x.closedTime == null &&
        x.props.savedFile == savedFile &&
        x.props.savedFolder == savedFolder &&
        x.props.savedFilePath == savedFilePath
    )
  }

  // @ts-ignore
  const {forceNewTab} = options || {};

  const component = tabs[newTab.tabComponent] //newTab.tabComponent TableDataTab
  if (!existing && !forceNewTab && component && component.matchingProps) {
    const testString = stableStringify(pick(newTab.props || {}, component.matchingProps))
    existing = unref(oldTabs).find(
      x =>
        x.props &&
        x.tabComponent == newTab.tabComponent &&
        x.closedTime == null &&
        stableStringify(pick(x.props || {}, component.matchingProps)) == testString
    )
  }

  if (existing) {
    // @ts-ignore
    locale.updateOpenedTabs(tabs => setSelectedTabFunc(tabs, existing.tabid))
    return
  }

  // new tab will be created
  if (newTab.title.endsWith('#')) {
    const numbers = unref(oldTabs)
      .filter(x => x.closedTime == null && x.title && x.title.startsWith(newTab.title))
      .map(x => parseInt(x.title.substring(newTab.title.length)));

    newTab.title = `${newTab.title}${findFreeNumber(numbers)}`;
  }

  const tabid = buildUUID()
  if (initialData) {
    for (const key of keys(initialData)) {
      if (key == 'editor') {
        await localforage.setItem(`tabdata_${key}_${tabid}`, initialData[key])
      } else {
        localStorage.setItem(`tabdata_${key}_${tabid}`, JSON.stringify(initialData[key]))
      }
    }
  }

  locale.updateOpenedTabs(files => {
    const dbKey = getTabDbKey(newTab)
    const items = sortTabs(files.filter(x => x.closedTime == null));

    const newItem = {
      ...newTab,
      tabid,
    }
    if (dbKey != null) {
      const lastIndex = findLastIndex(items, x => getTabDbKey(x) == dbKey)
      if (lastIndex >= 0) {
        items.splice(lastIndex + 1, 0, newItem);
      } else {
        items.push(newItem)
      }
    } else {
      items.push(newItem);
    }

    return [
      ...(files || []).map(x => ({
        ...x,
        selected: false,
        tabOrder: findIndex(items, y => y.tabid == x.tabid)
      })),
      {
        ...newTab,
        tabid,
        selected: true,
        tabOrder: findIndex(items, y => y.tabid == tabid),
      }
    ]
  })
}
