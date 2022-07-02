import stableStringify from 'fast-safe-stringify';
import {apiCall} from '/@/second/utility/api'
import {loadCachedValue} from './cache'

const connectionInfoLoader = ({conid}) => ({
  url: 'bridge.Connections.Get',
  params: {conid},
  reloadTrigger: 'connection-list-changed',
})

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

const serverStatusLoader = () => ({
  url: 'bridge.ServerConnections.ServerStatus',
  params: null,
  reloadTrigger: `server-status-changed`,
})

const databaseListLoader = ({ conid }) => ({
  url: 'bridge.ServerConnections.ListDatabases',
  params: conid,
  reloadTrigger: `database-list-changed-${conid}`,
})

async function getCore(loader, args) {
  const {url, params} = loader(args);
  const key = stableStringify({url, ...params})

  async function doLoad() {
    return await apiCall(url, params)
  }

  return await loadCachedValue(key, doLoad)
}

export function useCore(loader, args) {
  // const { url, params, reloadTrigger, transform, onLoaded } = loader(args);
  // const cacheKey = stableStringify({ url, ...params });
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

export function getConnectionList() {
  return getCore(connectionListLoader, {});
}

export function getConnectionInfo(args) {
  return getCore(connectionInfoLoader, args);
}

export function getServerStatus() {
  return getCore(serverStatusLoader, {})
}

export function useServerStatus() {
  return useCore(serverStatusLoader, {})
}

export function getDatabaseList(args) {
  return getCore(databaseListLoader, args);
}

export function useDatabaseList(args) {
  return useCore(databaseListLoader, args);
}
