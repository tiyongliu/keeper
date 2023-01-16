<template>
  <ConfigProvider :locale="getAntdLocale">
    <AppProvider>
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
// import 'dayjs/locale/zh-cn';

import 'ant-design-vue/es/input/style/css'
import 'ant-design-vue/es/modal/style/css'
import 'ant-design-vue/es/menu/style/css'
import 'ant-design-vue/es/tag/style/css'
import 'ant-design-vue/es/table/style/css'
import 'ant-design-vue/es/checkbox/style/css'
import 'ant-design-vue/es/select/style/css'
import 'ant-design-vue/es/form/style/css'
import 'ant-design-vue/es/row/style/css'
import 'ant-design-vue/es/col/style/css'
import 'ant-design-vue/es/radio/style/css'
import 'ant-design-vue/es/input-number/style/css'

// import 'ant-design-vue/es/icon/style/css'

let loadedApi = ref(false)

// support Multi-language
const {getAntdLocale} = useLocale();

initPluginsProvider()

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
