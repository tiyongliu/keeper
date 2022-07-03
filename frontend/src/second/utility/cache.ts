export async function loadCachedValue(cacheKey, func) {
  if (cacheKey) {}
  return await func()
}
