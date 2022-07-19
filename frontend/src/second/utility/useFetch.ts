import {onMounted, ref} from 'vue'
import stableStringify from 'json-stable-stringify';
import {isEqual} from 'lodash-es'

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

  const indicators = [url, stableStringify(data), stableStringify(params), loadCounter]

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

    const res = await doLoad()
    value.value = [res, loadedIndicators]
  }

  onMounted(() => {
    void loadValue(indicators)
  })

  const [returnValue, loadedIndicators] = value.value
  if (isEqual(indicators, loadedIndicators)) return returnValue

  return defaultValue
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
