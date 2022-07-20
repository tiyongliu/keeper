import stableStringify from 'json-stable-stringify'
import useFetch from '/@/second/utility/useFetch'

const connectionListLoader = () => ({
  url: 'bridge.Connections.List',
  params: {},
  reloadTrigger: `connection-list-changed`,
});

const serverStatusLoader = () => ({
  url: 'bridge.ServerConnections.ServerStatus',
  params: {},
  reloadTrigger: `server-status-changed`,
})

export function useConnectionList() {
  return useCore(connectionListLoader, {});
}

export function useServerStatus() {
  return useCore(serverStatusLoader, {});
}

function useCore(loader, args) {
  const { url, params, reloadTrigger, transform } = loader(args);
  const cacheKey = stableStringify({ url, ...params });

  const res = useFetch({
    url,
    params,
    reloadTrigger,
    cacheKey,
    transform,
  })

  return res
}

