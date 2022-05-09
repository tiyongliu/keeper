<template>
<!--  <a-tabs default-active-key="1" v-for="(tab, index) in _tabs">-->
<!--    <a-tab-pane :key="index + 1" :tab="tab.label"/>-->
<!--  </a-tabs>-->
</template>

<script lang="ts">
import {defineComponent, toRefs} from 'vue'
import {Tabs} from 'ant-design-vue'
import {compact} from 'lodash-es'

const TabPane = Tabs.TabPane

interface TabDef {
  label: string;
  slot?: number;
  component?: any;
  props?: any;
}

export default defineComponent({
  name: 'TabControl',
  components: {
    Tabs,
    TabPane
  },
  props: {
    isInline: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    tabs: {
      type: Array as PropType<TabDef[]>
    }
  },
  setup(props) {
    // const {tabs} = toRefs(props)

    console.log(`tabs.value`, props.tabs)
    return {
      _tabs: compact(props.tabs)
    }
  }
})
</script>

<style scoped>

</style>
