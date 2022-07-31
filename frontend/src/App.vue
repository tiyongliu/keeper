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
  import {onMounted, onBeforeUnmount} from 'vue'
  import {loadDatabasesApi} from '/@/api/connection'
  import {handleResetApi} from '/@/api/serverConnections'

  //TODO
  import { subscribeConnectionPingers } from '/@/api/connectionsPinger';
  let loadedApi = false

  import 'dayjs/locale/zh-cn';
  // support Multi-language
  const { getAntdLocale } = useLocale();

  // Listening to page changes and dynamically changing site titles
  useTitle()

  async function loadApi() {
    try {
      const connections = await loadDatabasesApi()
      if (connections) {
        loadedApi = true
      }

      if (loadedApi) {
        subscribeConnectionPingers()
      }

      if (!loadedApi) {
        console.log('API not initialized correctly, trying again in 1s');
        setTimeout(loadApi, 1000000);
      }
    } catch (err) {
      console.log('Error calling API, trying again in 1s', err);
      setTimeout(loadApi, 1000000);
    }
  }

  onMounted(() => {
    handleResetApi()
    loadApi()
    const removed = document.getElementById('starting_dbgate_zero');
    if (removed) removed.remove();
  })

</script>

<!--
https://vueuse.org/rxjs/useObservable/#usage

https://cn.rx.js.org/manual/tutorial.html#h22

https://pinia.web3doc.top/introduction.html#%E4%B8%80%E4%B8%AA%E6%9B%B4%E7%8E%B0%E5%AE%9E%E7%9A%84%E4%BE%8B%E5%AD%90
-->
