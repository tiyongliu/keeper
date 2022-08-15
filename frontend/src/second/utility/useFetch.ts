import {computed, onBeforeUnmount, ref, watch} from 'vue'
import stableStringify from 'json-stable-stringify';
import {isEqual} from 'lodash-es'
import {cacheClean, cacheGet, cacheSet, getCachedPromise} from '/@/second/utility/cache'
import getAsArray from '/@/second/utility/getAsArray'
import {useSocket} from '/@/second/utility/SocketProvider'
const apiLogging = true

export default function useFetch({
                                   url,
                                   data = undefined,
                                   params = undefined,
                                   defaultValue = undefined,
                                   reloadTrigger = undefined,
                                   cacheKey = undefined,
                                   transform = x => x,
                                   events,
                                   ...config
                                 }) {
  const value = ref([defaultValue, []])
  const loadCounter = ref(0)
  const socket = useSocket();
  const indicators = [url, stableStringify(data), stableStringify(params), loadCounter]

  async function loadValue(loadedIndicators) {
    async function doLoad() {
      if (apiLogging) {
        console.log('>>> API CALL', url, params)
      }
      try {
        let self: Function = window['go'];
        url.split(/[.\/]/).filter(item => item).forEach(key => self = self[key])
        if (!params || Object.keys(params).length === 0) {
          const resp = await self()
          return processApiResponse(url, params, resp)
        }
      } catch (e) {
        return Promise.reject(e)
      }
    }

    if (cacheKey) {
      const fromCache = cacheGet(cacheKey);
      if (fromCache) {
        value.value = [fromCache, loadedIndicators]
      } else {
        try {
          const res = await getCachedPromise(cacheKey, doLoad);
          cacheSet(cacheKey, res, reloadTrigger);
          value.value = [res, loadedIndicators]
        } catch (err) {
          console.error('Error when using cached promise', err);
          cacheClean(cacheKey);
          const res = await doLoad();
          cacheSet(cacheKey, res, reloadTrigger);
          value.value = [res, loadedIndicators]
        }
      }
    } else {
      const res = await doLoad()
      value.value = [res, loadedIndicators]
    }
  }

  watch(() => indicators, () => {
    void loadValue(indicators)
  }, {
    deep: true,
    immediate: true
  })

  watch(socket, () => {
    if (reloadTrigger && socket.value) {
      for (const item of getAsArray(reloadTrigger)) {
        socket.value.on(item, () => {
          void loadValue(indicators)
        })
      }
    }
  }, {
    immediate: true
  })


  onBeforeUnmount(() => {
    console.log(`onBeforeUnmount,onBeforeUnmount`, reloadTrigger)
    for (const item of getAsArray(reloadTrigger)) {
      if (socket.value) {
        socket.value.off(item)
      }
    }
  })

  return computed(() => {
    const [returnValue, loadedIndicators] = value.value
    if (isEqual(indicators, loadedIndicators)) return returnValue
    return defaultValue
  })
}

function processApiResponse(relativePath, params, resp) {
  if (apiLogging) {
    console.log('<<< API RESPONSE', relativePath, params, resp)
  }

  if (resp.status === 1) {
    // return resp.result.message
    throw resp.message
  }
  return resp.result
}
