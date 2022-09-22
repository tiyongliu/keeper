<template>
  <ConfigProvider :locale="getAntdLocale">
    <AppProvider>
      <RouterView />
    </AppProvider>
  </ConfigProvider>
</template>

<script lang="ts" setup>
  import { ConfigProvider } from 'ant-design-vue';
  import { AppProvider } from '/@/components/Application';
  import { useTitle } from '/@/hooks/web/useTitle';
  import { useLocale } from '/@/locales/useLocale';
  import {onMounted} from 'vue'
  import {connectionListApi} from '/@/api/simpleApis'

  //TODO
  import initPluginsProvider from '/@/second/plugins/PluginsProvider'
  import { subscribeConnectionPingers } from '/@/api/connectionsPinger';
  let loadedApi = false

  import 'dayjs/locale/zh-cn';
  // support Multi-language
  const { getAntdLocale } = useLocale();

  // Listening to page changes and dynamically changing site titles
  useTitle()

  initPluginsProvider()

  async function loadApi() {
    try {
      const connections = await connectionListApi()
      if (connections) {
        loadedApi = true
      }

      if (loadedApi) {
        subscribeConnectionPingers()
      }

      if (!loadedApi) {
        console.log('API not initialized correctly, trying again in 1s');
        setTimeout(loadApi, 1000);
      }
    } catch (err) {
      console.log('Error calling API, trying again in 1s', err);
      setTimeout(loadApi, 1000);
    }
  }

  onMounted(() => {
    loadApi()
    const removed = document.getElementById('starting_dbgate_zero');
    if (removed) removed.remove();
  })

</script>
