import getAsArray from './getAsArray';

const cachedByKey = {};
const cachedPromisesByKey = {};
const cachedKeysByReloadTrigger = {};
const subscriptionsByReloadTrigger = {};
const cacheGenerationByKey = {};

let cacheGeneration = 0;

export function cacheGet(key) {
  return cachedByKey[key];
}

export function cacheSet(key, value, reloadTrigger) {
  cachedByKey[key] = value;
  for (const item of getAsArray(reloadTrigger)) {
    if (!(item in cachedKeysByReloadTrigger)) {
      cachedKeysByReloadTrigger[item] = [];
    }
    cachedKeysByReloadTrigger[item].push(key);
  }
  delete cachedPromisesByKey[key];
}

export function cacheClean(reloadTrigger) {
  for (const item of getAsArray(reloadTrigger)) {
    const keys = cachedKeysByReloadTrigger[item];
    if (keys) {
      for (const key of keys) {
        delete cachedByKey[key];
        delete cachedPromisesByKey[key];
      }
    }
    delete cachedKeysByReloadTrigger[item];
  }
}


function addCacheKeyToReloadTrigger(cacheKey, reloadTrigger) {
  for (const item of getAsArray(reloadTrigger)) {
    if (!(item in cachedKeysByReloadTrigger)) {
      cachedKeysByReloadTrigger[item] = [];
    }
    cachedKeysByReloadTrigger[item].push(cacheKey);
  }
}

function acquireCacheGeneration() {
  cacheGeneration += 1;
  return cacheGeneration;
}

function getCacheGenerationForKey(cacheKey) {
  return cacheGenerationByKey[cacheKey] || 0;
}

export async function loadCachedValue(reloadTrigger, cacheKey, func) {
  // const fromCache = cacheGet(cacheKey)
  // if (fromCache) {
  //   return fromCache;
  // } else {
  //   const generation = acquireCacheGeneration();
  //   try {
  //     const res = await getCachedPromise(reloadTrigger, cacheKey, func);
  //     if (getCacheGenerationForKey(cacheKey) > generation) {
  //       return cacheGet(cacheKey) || res;
  //     } else {
  //       cacheSet(cacheKey, res, reloadTrigger, generation);
  //       return res;
  //     }
  //   } catch (err) {
  //     console.error('Error when using cached promise', err);
  //     cacheClean(cacheKey);
  //     const res = await func();
  //     cacheSet(cacheKey, res, reloadTrigger, generation);
  //     return res;
  //   }
  // }


  if (cacheKey) {}
  return await func()
}

export async function subscribeCacheChange(reloadTrigger, cacheKey, reloadHandler) {
  for (const item of getAsArray(reloadTrigger)) {
    if (!subscriptionsByReloadTrigger[item]) {
      subscriptionsByReloadTrigger[item] = [];
    }
    subscriptionsByReloadTrigger[item].push(reloadHandler);
  }
}

export async function unsubscribeCacheChange(reloadTrigger, cacheKey, reloadHandler) {
  for (const item of getAsArray(reloadTrigger)) {
    if (subscriptionsByReloadTrigger[item]) {
      subscriptionsByReloadTrigger[item] = subscriptionsByReloadTrigger[item].filter(x => x != reloadHandler);
    }
    if (subscriptionsByReloadTrigger[item].length == 0) {
      delete subscriptionsByReloadTrigger[item];
    }
  }
}

export function getCachedPromise(key, func) {
  if (key in cachedPromisesByKey) return cachedPromisesByKey[key];
  const promise = func();
  cachedPromisesByKey[key] = promise;
  return promise;
}
