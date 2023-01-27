<template>
  <ConfigProvider :locale="getAntdLocale">
    <AppProvider>
      <DataGridRowHeightMeter />
      <RouterView />
    </AppProvider>
  </ConfigProvider>
</template>

<script lang="ts" setup>
import {ConfigProvider} from 'ant-design-vue';
import {AppProvider} from '/@/components/Application';
import {useLocale} from '/@/locales/useLocale';
import {onMounted, ref, watchEffect} from 'vue'
import {storeToRefs} from 'pinia'
import {connectionListApi} from '/@/api/simpleApis'
import dispatchRuntimeEvent from "/@/api/event";
import initPluginsProvider from '/@/second/plugins/PluginsProvider'
import {subscribeConnectionPingers} from '/@/api/connectionsPinger'
import {setAppLoaded} from '/@/second/utility/appLoadManager'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {useAppStore} from '/@/store/modules/app'
import DataGridRowHeightMeter from '/@/second/datagrid/DataGridRowHeightMeter.vue'
// import 'dayjs/locale/zh-cn';

let loadedApi = ref(false)

// support Multi-language
const {getAntdLocale} = useLocale();

initPluginsProvider()

const appStore = useAppStore()
const bootstrap = useBootstrapStore()
const {loadingPluginStore} = storeToRefs(bootstrap)

async function loadApi() {
  try {
    const connections = await connectionListApi()
    if (connections) {
      loadedApi.value = true
    }

    if (loadedApi.value) {
      subscribeConnectionPingers()
    }

    if (!loadedApi.value) {
      console.log('API not initialized correctly, trying again in 1s');
      setTimeout(loadApi, 1000);
    }
  } catch (err) {
    console.log('Error calling API, trying again in 1s', err);
    setTimeout(loadApi, 1000);
  }
}

watchEffect(() => {
  if (loadedApi.value && loadingPluginStore.value.loaded) {
    setAppLoaded();
  }
})

onMounted(() => {
  loadApi()
  const removed = document.getElementById('starting_dbgate_zero');
  if (removed) removed.remove();

  if (window.runtime) {
    dispatchRuntimeEvent()
  }
})
</script>
