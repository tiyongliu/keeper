import {computed, ref, watch, watchEffect} from 'vue'
import stableStringify from 'json-stable-stringify';
import {isEqual} from 'lodash-es'
import {EventsOn} from '/@/wailsjs/runtime/runtime'
import getAsArray from '/@/second/utility/getAsArray'
import {cacheClean, cacheGet, cacheSet, getCachedPromise} from '/@/second/utility/cache'

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

  const indicators = [url, stableStringify(data), stableStringify(params), loadCounter.value]

  async function loadValue(loadedIndicators) {
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

  watch(reloadTrigger, () => {
    if (reloadTrigger && !events) {
      console.error('Socket not available, reloadTrigger not planned')
    }
    if (reloadTrigger && events) {
      const {eventsOn} = events
      if (eventsOn) {
        EventsOn(reloadTrigger, async (data) => {
          console.log(`data111111111111111111`, data)
        })
      }
    }
  }, {
    immediate: true
  })

  const socketEvent = () => {

  }
  socketEvent()

  watch(() => indicators, () => {
    void loadValue(indicators)
  }, {
    // deep: true,
    immediate: true
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
