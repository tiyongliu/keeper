import type { App } from 'vue';
import { createPinia } from 'pinia';
import {safeJsonParse} from "/@/utils/lib/pkg/stringTools";

const store = createPinia();

export function setupStore(app: App<Element>) {
  app.use(store);
}

export { store };

export function getWithStorageVariableCache<T>(defaultValue: T, storageName) {
  const init = localStorage.getItem(storageName);
  return (init ? safeJsonParse(init, defaultValue, true) : defaultValue);
}

export function setWithStorageVariableCache(storageName, value: any) {
  localStorage.setItem(storageName, JSON.stringify(value));
}
