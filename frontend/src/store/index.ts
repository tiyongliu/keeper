import type {App} from 'vue';
import {createPinia} from 'pinia';

const store = createPinia();

export function setupStore(app: App<Element>) {
  app.use(store);
}

export {store}

const currentSettingsValue = {
  app: {
    useNativeMenu: false
  }
}
//todo 临时写死
export const getCurrentSettings = () => currentSettingsValue || {};
