import {ref, Ref} from 'vue'
import stableStringify from 'json-stable-stringify'
import {extendDatabaseInfo} from '/@/second/keeper-tools'
import {setLocalStorage} from '/@/second/utility/storageCache'
import {EventsOn} from '/@/wailsjs/runtime/runtime'
import getAsArray from '/@/second/utility/getAsArray'
import {apiCall} from '/@/second/utility/api'
import {loadCachedValue} from '/@/second/utility/cache'

const connectionListLoader = () => ({
  url: 'bridge.Connections.List',
  params: {},
  reloadTrigger: `connection-list-changed`
});

const serverStatusLoader = () => ({
  url: 'bridge.ServerConnections.ServerStatus',
  params: {},
  reloadTrigger: `server-status-changed`,
})

const databaseListLoader = ({conid}) => ({
  url: 'bridge.ServerConnections.ListDatabases',
  params: {conid},
  reloadTrigger: `database-list-changed-${conid}`,
  onLoaded: value => {
    if (value?.length > 0) setLocalStorage(`database_list_${conid}`, value);
  },
})

const databaseServerVersionLoader = ({conid, database}) => ({
  url: 'bridge.DatabaseConnections.ServerVersion',
  params: {conid, database},
  reloadTrigger: `database-server-version-changed-${conid}-${database}`,
})

const databaseStatusLoader = ({conid, database}) => ({
  url: 'bridge.DatabaseConnections.Status',
  params: {conid, database},
  reloadTrigger: `database-status-changed-${conid}-${database}`,
})

const databaseInfoLoader = ({conid, database}) => ({
  url: 'bridge.DatabaseConnections.Structure',
  params: {conid, database},
  reloadTrigger: `database-structure-changed-${conid}-${database}`,
  transform: extendDatabaseInfo,
})

const connectionInfoLoader = ({conid}) => ({
  url: 'bridge.Connections.Get',
  params: {conid},
  reloadTrigger: 'connection-list-changed',
})

const installedPluginsLoader = () => ({
  url: 'bridge.Plugins.Installed',
  params: {},
  reloadTrigger: `installed-plugins-changed`,
})

export function useConnectionList<T>(targetRef: Ref<T>) {
  return useCore(connectionListLoader, {}, targetRef);
}

export function useServerStatus<T>(targetRef: Ref<T>) {
  return useCore(serverStatusLoader, {}, targetRef);
}

export function useDatabaseList<T>(args, targetRef: Ref<T>) {
  return useCore(databaseListLoader, args, targetRef);
}

export function useDatabaseServerVersion<T>(args, targetRef: Ref<T>) {
  return useCore(databaseServerVersionLoader, args, targetRef);
}

export function useDatabaseStatus<T>(args, targetRef: Ref<T>) {
  return useCore(databaseStatusLoader, args, targetRef);
}

export function useDatabaseInfo<T>(args, targetRef: Ref<T>) {
  return useCore(databaseInfoLoader, args, targetRef);
}

export function useConnectionInfo<T>(args, targetRef: Ref<T>) {
  return useCore(connectionInfoLoader, args, targetRef)
}

export function useInstalledPlugins<T>(args = {},targetRef: Ref<T>) {
  return useCore(installedPluginsLoader, args, targetRef);
}

async function getCore(loader, args) {
  const {url, params, reloadTrigger, transform, onLoaded, errorValue} = loader(args);
  const key = stableStringify({url, ...params});

  async function doLoad() {
    const resp = await apiCall(url, params);
    if (resp?.errorMessage && errorValue !== undefined) {
      if (onLoaded) onLoaded(errorValue)
      return errorValue;
    }
    const res = (transform || (x => x))(resp)
    if (onLoaded) onLoaded(res);
    return res
  }

  return await loadCachedValue(reloadTrigger, key, doLoad)
}

function useCore<T>(loader, args, targetRef: Ref<T | null | undefined>) {
  const openedCount = ref(0)
  const {reloadTrigger} = loader(args);

  async function handleReload() {
    const res = await getCore(loader, args);
    if (openedCount.value > 0) {
      targetRef.value = res
    }
  }

  if (reloadTrigger) {
    for (const item of getAsArray(reloadTrigger)) {
      try {
        EventsOn(item, () => {
          void handleReload()
        })
      } catch (e) {
        console.log(e)
      }
    }
  }

  openedCount.value += 1
  void handleReload()
}
