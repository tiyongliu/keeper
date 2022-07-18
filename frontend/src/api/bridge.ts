import stableStringify from 'json-stable-stringify'
import {
  loadCachedValue,
  subscribeCacheChange,
  unsubscribeCacheChange
} from '/@/second/utility/cache'
import {apiCall} from '/@/second/utility/api'

const databaseListLoader = ({conid}) => ({
  url: 'bridge.ServerConnections.ListDatabases',
  params: {conid},
  reloadTrigger: `connection-list-changed-${conid}`
})

export function useDatabaseList(args) {
  return useCore(databaseListLoader, args)
}

async function getCore(loader, args) {
  const {url, params, reloadTrigger, transform, onLoaded, errorValue} = loader(args);
  const key = stableStringify({url, ...params});

  async function doLoad() {
    return await apiCall(url, params)
  }

  const res = await loadCachedValue(reloadTrigger, key, doLoad)
  return res
}

function useCore(loader, args) {
  return new Promise(resolve => {
    const {url, params, reloadTrigger, transform, onLoaded} = loader(args);
    const cacheKey = stableStringify({url, ...params})
    let openedCount = 0

    async function handleReload() {
      const res = await getCore(loader, args);
      if (openedCount > 0) {
        resolve(res)
      }
    }

    openedCount += 1
    void handleReload()

    if (reloadTrigger) {
      void subscribeCacheChange(reloadTrigger, cacheKey, handleReload)

      return () => {
        openedCount -= 1
        void unsubscribeCacheChange(reloadTrigger, cacheKey, handleReload)
      }
    } else {
      return () => {
        openedCount -= 1
      }
    }
  })


}
