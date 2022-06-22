export async function loadCachedValue(cacheKey, func) {
  return await func()
}
