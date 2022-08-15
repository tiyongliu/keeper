import stableStringify from 'json-stable-stringify'
import useFetch from '/@/second/utility/useFetch'

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

const databaseListLoader = ({ conid }) => ({
  url: 'bridge.ServerConnections.ListDatabases',
  params: { conid },
  reloadTrigger: `database-list-changed-${conid}`,
})

export function useConnectionList() {
  return useCore(connectionListLoader, {});
}

export function useServerStatus() {
  return useCore(serverStatusLoader, {});
}

export function useDatabaseList(args) {
  return useCore(databaseListLoader, args);
}

function useCore(loader, args) {
  const {url, params, reloadTrigger, transform, events} = loader(args);
  const cacheKey = stableStringify({url, ...params});

  const res = useFetch({
    url,
    params,
    reloadTrigger,
    cacheKey,
    transform,
    events,
  })

  return res
}

