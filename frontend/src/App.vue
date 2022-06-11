<template>
  <ConfigProvider :locale="getAntdLocale">
    <AppProvider>
      <RouterView/>
    </AppProvider>
  </ConfigProvider>
</template>

<script lang="ts" setup>
import {ConfigProvider} from 'ant-design-vue';
import {AppProvider} from '/@/components/Application';
import {useTitle} from '/@/hooks/web/useTitle';
import {useLocale} from '/@/locales/useLocale';
import {onMounted} from 'vue'
import {loadDatabasesApi} from '/@/api/connection'
import 'dayjs/locale/zh-cn';
// support Multi-language
const {getAntdLocale} = useLocale();

// Listening to page changes and dynamically changing site titles
useTitle();


async function loadApi() {
  try {
    const connections = await loadDatabasesApi()
    console.log(connections, `loadDatabasesApi`)
  } catch(e) {
    console.log('Error calling API, trying again in 1s', e);
  }
}

onMounted(() => {
  // loadApi()
})
</script>

<!--
https://vueuse.org/rxjs/useObservable/#usage

https://cn.rx.js.org/manual/tutorial.html#h22

https://pinia.web3doc.top/introduction.html#%E4%B8%80%E4%B8%AA%E6%9B%B4%E7%8E%B0%E5%AE%9E%E7%9A%84%E4%BE%8B%E5%AD%90
-->
