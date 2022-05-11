import stableStringify from 'fast-safe-stringify';

const connectionListLoader = () =>({
  url: 'connections/list',
  params: {},
  reloadTrigger: `connection-list-changed`
})




export function useConnectionList(){
  // return useCore(connectionListLoader,{})

  return useCore(connectionListLoader,{})
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
