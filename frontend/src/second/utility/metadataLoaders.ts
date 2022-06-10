import stableStringify from 'fast-safe-stringify';
import {apiCall} from '/@/second/utility/api'
import {loadCachedValue} from './cache'

const connectionListLoader = () =>({
  url: 'connections/list',
  params: null,
  reloadTrigger: `connection-list-changed`
})

const databaseServerVersionLoader = ({ conid, database }) => ({
  url: 'database-connections/server-version',
  params: { conid, database },
  reloadTrigger: `database-server-version-changed-${conid}-${database}`,
});
const databaseStatusLoader = ({ conid, database }) => ({
  url: 'database-connections/status',
  params: { conid, database },
  reloadTrigger: `database-status-changed-${conid}-${database}`,
});

async function getCore(loader, args){
  const { url, params } = loader(args);
  const key = stableStringify({url, ...params})

  async function doload(){
    const resp = await apiCall(url, params);
    return resp
  }

  const res = await loadCachedValue(key, doload)
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

function useCore(loader, args){
  const { url, params, reloadTrigger, transform, onLoaded } = loader(args);
  const cacheKey = stableStringify({ url, ...params });
  let closed = false

  return{
    subscribe: onChange =>{
      async function handleReload(){
        const res = await getCore(loader,args);
        if(!closed){
          onChange(res);
        }
      }
      handleReload()
    }
  }
}

export function useConnectionList() {
  return useCore(connectionListLoader, {})
}
