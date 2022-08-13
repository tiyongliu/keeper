// @ts-ignore
import {computed, ComputedRef, onBeforeUnmount, ref, UnwrapRefSimple} from "vue";
import stableStringify from 'json-stable-stringify'
import {isEqual} from "lodash-es";
import {extendDatabaseInfo} from 'keeper-tools';
import {setLocalStorage} from '/@/second/utility/storageCache'
import {EventsOff, EventsOn} from '/@/wailsjs/runtime/runtime'
import getAsArray from '/@/second/utility/getAsArray'
import {apiCall} from '/@/second/utility/api'
import {loadCachedValue} from './cache'

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
  url: 'database-connections/server-version',
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

export function useConnectionList<T>(): ComputedRef<T> {
  return useCore(connectionListLoader, {});
}

export function useServerStatus<T>(): ComputedRef<T> {
  return useCore(serverStatusLoader, {});
}

export function useDatabaseList<T>(args): ComputedRef<T> {
  return useCore(databaseListLoader, args);
}

export function useDatabaseServerVersion(args) {
  return useCore(databaseServerVersionLoader, args);
}

export function useDatabaseStatus(args) {
  return useCore(databaseStatusLoader, args);
}

export function useDatabaseInfo(args) {
  return useCore(databaseInfoLoader, args);
}

async function getCore(loader, args) {
  const {url, params, reloadTrigger, transform, onLoaded, errorValue} = loader(args);
  const key = stableStringify({url, ...params});

  async function doLoad() {
   try {
     const resp = await apiCall(url, params);
     if (resp?.errorMessage && errorValue !== undefined) {
       if (onLoaded) onLoaded(errorValue)
       return errorValue;
     }
     if (url == 'bridge.DatabaseConnections.Structure') {
       console.log(resp, `333333333`, url)
     }
     const res = (transform || (x => x))(resp);
     if (onLoaded) onLoaded(res);
     return res;
   } catch (e) {
     console.log(e)
   }
  }

  return await loadCachedValue(reloadTrigger, key, doLoad)
}


function useCore<T>(loader, args): ComputedRef<UnwrapRefSimple<T> | null | undefined> {
  const value = ref<[T | null, any]>([null, []])
  const openedCount = ref(0)
  const {url, params, reloadTrigger} = loader(args);
  const cacheKey = stableStringify({url, ...params})
  const indicators = [url, cacheKey, stableStringify(params), openedCount]

  async function handleReload(loadedIndicators) {
    const res = await getCore(loader, args);
    if (openedCount.value > 0) {
      value.value = [res, loadedIndicators]
    }
  }

  if (reloadTrigger) {
    for (const item of getAsArray(reloadTrigger)) {
      try {
        EventsOn(item, () => {
          void handleReload(indicators)
        })
      } catch (e) {
        console.log(e)
      }
    }
  }

  openedCount.value += 1
  void handleReload(indicators)

  onBeforeUnmount(() => {
    value.value = [null, []]
    openedCount.value -= 1
    if (reloadTrigger) {
      EventsOff(reloadTrigger)
    }
  })

  return computed(() => {
    const [returnValue, loadedIndicators] = value.value
    if (isEqual(indicators, loadedIndicators)) return returnValue
    return undefined
  })
}
