import {safeJsonParse} from '/@/second/keeper-tools'

export function getWithStorageVariableCache<T>(defaultValue: T, storageName) {
  const init = localStorage.getItem(storageName);
  const res = (init ? safeJsonParse(init, defaultValue, true) : defaultValue)
  localStorage.setItem(storageName, JSON.stringify(res))
  return res
}

export function setWithStorageVariableCache(storageName, value: any) {
  localStorage.setItem(storageName, JSON.stringify(value));
}
