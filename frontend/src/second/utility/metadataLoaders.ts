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
    const resp = await apiCall(url, params);
    return resp
  }
  // const res = await doLoad()
  const res = await loadCachedValue(key, doLoad)
  console.log(`line 32 `, res)
  return res
}

// export function useConnectionList(){
//   // return useCore(connectionListLoader,{})
//
//   return useCore(connectionListLoader,{})
// }

export function useDatabaseServerVersion(args) {
  return useCore(databaseServerVersionLoader, args);
}

export function useDatabaseStatus(args) {
  return useCore(databaseStatusLoader, args);
}

function useCore(loader, args) {
  const {url, params, reloadTrigger, transform, onLoaded} = loader(args);
  console.log(url, params, `line 51-51`)
  const cacheKey = stableStringify({url, ...params});
  console.log(cacheKey, `cacheKey-cacheKey`)

  let closed = false

  return {
    subscribe: onChange => {
      async function handleReload() {

        console.log(loader, args, `rrrrrrrrrrrrrrrrr`)
        const res = await getCore(loader, args);

        if (!closed) {
          console.log(res, `res`)
          console.log(onChange, `oooo`)

          // onChange(res);
        }
      }

      void handleReload()
    }
  }
}

export function useConnectionList() {
  return useCore(connectionListLoader, {})
}

import { defineStore } from 'pinia'
