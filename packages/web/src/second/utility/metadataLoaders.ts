import stableStringify from 'fast-safe-stringify';

const connectionListLoader = () =>({
  url: 'connections/list',
  params: {},
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



export function useConnectionList(){
  // return useCore(connectionListLoader,{})

  return useCore(connectionListLoader,{})
}
export function useDatabaseServerVersion(args) {
  return useCore(databaseServerVersionLoader, args);
}
export function useDatabaseStatus(args) {
  return useCore(databaseStatusLoader, args);
}


function useCore(loader, args){
  const { url, params, reloadTrigger, transform, onLoaded } = loader(args);
  const cacheKey = stableStringify({ url, ...params });
  let closed = false;

  return{
    subscribe: onChange =>{
      async function handleReload(){
        const res = await getCore(loader,args);
        if(!closed){
          onChange(res);
        }
      }
      handleReload();


    }
  }

}

async function getCore(loader, args){
  const { url, params, reloadTrigger, transform, onLoaded, errorValue } = loader(args);
  const key = stableStringify({url, ...params})

  async function doload(){
    const resp = await apiCall(url,params);

  }

}
