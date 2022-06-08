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
import {apiCall} from '/@/second/utility/api'
import 'dayjs/locale/zh-cn';
// support Multi-language
const {getAntdLocale} = useLocale();

// Listening to page changes and dynamically changing site titles
useTitle();


async function loadApi() {
  try {
    const connections = await apiCall('bridge.Connections.List')
    // const connections = await window['go'].bridge.Connections.List()
    console.log(connections, `erterterterter`)
  } catch(e) {
    console.log('Error calling API, trying again in 1s', e);
  }
}

onMounted(() => {
  loadApi()
})
</script>
