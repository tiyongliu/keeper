import {onMounted, ref, watch, computed} from 'vue'
import stableStringify from 'json-stable-stringify';
import {isEqual} from 'lodash-es'
import {cacheClean, cacheGet, cacheSet, getCachedPromise} from '/@/second/utility/cache'
import {useRefs} from '/@/hooks/core/useRefs'
const apiLogging = true

export default function useFetch({
                                   url,
                                   data = undefined,
                                   params = undefined,
                                   defaultValue = undefined,
                                   reloadTrigger = undefined,
                                   cacheKey = undefined,
                                   transform = x => x,
                                   ...config
                                 }) {
  const value = ref([defaultValue, []])
  const loadCounter = ref(0)


  const indicators = [url, stableStringify(data), stableStringify(params), loadCounter.value]

  async function loadValue(loadedIndicators) {
    console.log(`111`)
    async function doLoad() {
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
      console.log(`222`)
      const fromCache = cacheGet(cacheKey);
      if (fromCache) {
        value.value = [fromCache, loadedIndicators]
        console.log(`444`, value.value)
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
      console.log(`333`)
      const res = await doLoad()
      value.value = [res, loadedIndicators]
      console.log(res, `[res, loadedIndicators]`)
    }
  }

  const socketEvent = () => {
    if (reloadTrigger) {
      console.error('Socket not available, reloadTrigger not planned')
    }
  }

  onMounted(() => {
    void loadValue(indicators)
  })

  watch(() => indicators, () => {
    console.log(`ffff`)
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
