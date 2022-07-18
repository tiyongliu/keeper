import {ref} from 'vue'
import stableStringify from 'json-stable-stringify';

export default function useFetch({
 url,
 data = undefined,
 params = undefined,
 defaultValue = undefined,
 reloadTrigger = undefined,
 cacheKey = undefined,
 transform = x => x,
 ...config
})  {
  const value = ref([defaultValue, []])
  const loadCounter = ref(0)

  const indicators =  [url, stableStringify(data), stableStringify(params), loadCounter]

  async function loadValue(loadedIndicators) {

  }

}
