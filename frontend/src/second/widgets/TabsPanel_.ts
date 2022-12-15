import {nextTick, unref} from 'vue'
import {storeToRefs} from "pinia"
import {get, findLastIndex} from 'lodash-es'
import {useLocaleStore} from "/@/store/modules/locale"
import getConnectionLabel from '/@/second/utility/getConnectionLabel'

const locale = useLocaleStore()
const {getOpenedTabs} = storeToRefs(locale)

const closeTabFunc = closeCondition => tabid => {
  const files = unref(getOpenedTabs)
  locale.setOpenedTabs(() => {
    const active = files.find(x => x.tabid == tabid);
    if (!active) return files;
    const newFiles = files.map(x => ({
      ...x,
      closedTime: x.closedTime || (closeCondition(x, active) ? new Date().getTime() : undefined),
    }));

    if (newFiles.find(x => x.selected && x.closedTime == null)) {
      return newFiles;
    }

    const selectedIndex = findLastIndex(newFiles, x => x.closedTime == null)
    return newFiles.map((x, index) => ({
      ...x,
      selected: index == selectedIndex,
    }))
  })
}

export const closeMultipleTabs = (closeCondition, deleteFromHistory = false) => {
  const files = unref(getOpenedTabs)
  locale.setOpenedTabs(() => {
    const newFiles = deleteFromHistory
      ? files.filter(x => !closeCondition(x))
      : files.map(x => ({
        ...x,
        closedTime: x.closedTime || (closeCondition(x) ? new Date().getTime() : undefined),
      }));

    if (newFiles.find(x => x.selected && x.closedTime == null)) {
      return newFiles;
    }

    const selectedIndex = findLastIndex(newFiles, x => x.closedTime == null)
    return newFiles.map((x, index) => ({
      ...x,
      selected: index == selectedIndex,
    }))
  })
}

export const closeTab = closeTabFunc((x, active) => x.tabid == active.tabid);

export const closeWithSameDb = closeTabFunc(
  (x, active) =>
    get(x, 'props.conid') == get(active, 'props.conid') &&
    get(x, 'props.database') == get(active, 'props.database')
)

export const closeWithOtherDb = closeTabFunc(
  (x, active) =>
    get(x, 'props.conid') != get(active, 'props.conid') ||
    get(x, 'props.database') != get(active, 'props.database')
)

export function getTabDbName(tab, connectionList) {
  if (tab.tabComponent == 'ConnectionTab') return 'Connections';
  if (tab.props && tab.props.conid && tab.props.database) return tab.props.database;
  if (tab.props && tab.props.conid) {
    const connection = connectionList?.find(x => x._id == tab.props.conid);
    if (connection) return getConnectionLabel(connection, { allowExplicitDatabase: false });
    return '???';
  }
  if (tab.props && tab.props.archiveFolder) return tab.props.archiveFolder;
  return '(no DB)';
}

export async function scrollInViewTab(tabid) {
  await nextTick()
  const element = document.getElementById(`file-tab-item-${tabid}`);
  if (element) {
    element.scrollIntoView({ block: 'nearest', inline: 'nearest' });
  }
}

export function getDbIcon(key) {
  if (key) {
    if (key.startsWith('database://')) return 'icon database';
    if (key.startsWith('archive://')) return 'icon archive';
    if (key.startsWith('server://')) return 'icon server';
    if (key.startsWith('connections.')) return 'icon connection';
  }
  return 'icon file';
}
