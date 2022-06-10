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
