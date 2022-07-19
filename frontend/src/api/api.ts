import stableStringify from 'json-stable-stringify'
import useFetch from '/@/second/utility/useFetch'

const connectionListLoader = () => ({
  url: 'bridge.Connections.List',
  params: {},
  reloadTrigger: `connection-list-changed`,
});

export function useConnectionList() {
  return useCore(connectionListLoader, {});
}

function useCore(loader, args) {
  const { url, params, reloadTrigger, transform } = loader(args);
  const cacheKey = stableStringify({ url, ...params })

  const res = useFetch({
    url,
    params,
    reloadTrigger,
    cacheKey,
    transform,
  })

  return res
}

