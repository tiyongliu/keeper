<template>
  <div class="wrapper">
    <div class="section">
      <WidgetTitle>Execute</WidgetTitle>
    </div>

    <div class="section">
      <WidgetTitle>Parameters</WidgetTitle>
      <template v-if="selectedMacro?.args && selectedMacro?.args?.length > 0">

      </template>
      <div class="m-1" v-else>This macro has no parameters</div>
    </div>

    <div class="section">
      <WidgetTitle>Description</WidgetTitle>
      <div class="m-1">{{selectedMacro?.description}}}</div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, inject} from 'vue'
import WidgetTitle from '/@/second/widgets/WidgetTitle.vue'
export default defineComponent({
  name: "MacroInfoTab",
  props: {
    onExecute: {
      type: Function as PropType<() => void>
    }
  },
  components: {
    WidgetTitle
  },
  setup(props) {
    const selectedMacro = inject('selectedMacro') as any;
    return {
      ...props,
      selectedMacro,
    }
  }
})
</script>

<style scoped>
.wrapper {
  display: flex;
  overflow-y: auto;
}
.section {
  margin: 5px;
}
</style>
