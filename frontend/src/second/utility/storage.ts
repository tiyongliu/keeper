import {safeJsonParse} from '/@/second/keeper-tools'

export function getWithStorageVariableCache<T>(defaultValue: T, storageName) {
  const init = localStorage.getItem(storageName);
  return (init ? safeJsonParse(init, defaultValue, true) : defaultValue);
}

export function setWithStorageVariableCache(storageName, value: any) {
  localStorage.setItem(storageName, JSON.stringify(value));
}
