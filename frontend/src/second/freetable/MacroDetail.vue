<template>
  <div class="container">
    <MacroHeader :onExecute="onExecute"/>
    <TabControl :tabs="tabs"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, toRefs} from 'vue'
import MacroHeader from './MacroHeader.vue'
import TabControl from '/@/second/elements/TabControl.vue'
export default defineComponent({
  name: "MacroDetail",
  components: {
    MacroHeader,
    TabControl
  },
  props: {
    onExecute: {
      type: Function as PropType<() => void>
    }
  },
  setup(props) {
    return {
      ...toRefs(props),
      tabs: []
    }
  }
})
</script>

<style scoped>
.container {
  position: absolute;
  display: flex;
  flex-direction: column;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>
