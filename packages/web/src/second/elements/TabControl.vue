<template>
  <div>
    <div>
      <a-tabs default-active-key="1">
        <a-tab-pane :key="`${index + 1}`" :tab="tab.label" v-for="(tab, index) in tabs">
          <component :is="tab.component"/>
        </a-tab-pane>
      </a-tabs>
    </div>

    <div></div>
  </div>
</template>

<script lang="ts">
import {Component, defineComponent} from "vue";
import {compact} from "lodash-es";
import { Tabs} from "ant-design-vue";
const TabPane = Tabs.TabPane
interface TabDef {
  label: string;
  slot?: number;
  component?: string | Component;
  props?: any;
}
export default defineComponent({
  name: 'TabControl',
  components: {
    [Tabs.name]: Tabs,
    [TabPane.name]: TabPane,
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
    return {
      tabs: compact(props.tabs)
    }
  }
})
</script>

<style scoped>

</style>
