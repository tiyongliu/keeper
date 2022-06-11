import stableStringify from 'fast-safe-stringify';
import {apiCall} from '/@/second/utility/api'
import {loadCachedValue} from './cache'

const connectionListLoader = () => ({
  url: 'bridge.Connections.List',
  params: null,
  reloadTrigger: `connection-list-changed`
})

const databaseServerVersionLoader = ({conid, database}) => ({
  url: 'database-connections/server-version',
  params: {conid, database},
  reloadTrigger: `database-server-version-changed-${conid}-${database}`,
});
const databaseStatusLoader = ({conid, database}) => ({
  url: 'database-connections/status',
  params: {conid, database},
  reloadTrigger: `database-status-changed-${conid}-${database}`,
});

async function getCore(loader, args) {
  const {url, params} = loader(args);
  const key = stableStringify({url, ...params})

  async function doLoad() {
    return await apiCall(url, params)
  }

  return await loadCachedValue(key, doLoad)
}

export function useCore(loader, args) {
  const { url, params, reloadTrigger, transform, onLoaded } = loader(args);
  const cacheKey = stableStringify({ url, ...params });
  let closed = false;
  async function handleReload() {
    const res = await getCore(loader, args);
    if (!closed) {
      return res
    }
  }
  return handleReload()
}

export function useDatabaseServerVersion(args) {
  return useCore(databaseServerVersionLoader, args);
}

export function useDatabaseStatus(args) {
  return useCore(databaseStatusLoader, args);
}

export function useConnectionList() {
  return useCore(connectionListLoader, {})
}




export function useConnectionList1() {
  return useCore1(connectionListLoader, {})
}

export function useCore1(loader, args) {
  const { url, params, reloadTrigger, transform, onLoaded } = loader(args);
  const cacheKey = stableStringify({ url, ...params });
  let closed = false
  async function handleReload() {
    const res = await getCore(loader, args);
    if (!closed) {
      return res
    }
  }
  return {

  }
}
