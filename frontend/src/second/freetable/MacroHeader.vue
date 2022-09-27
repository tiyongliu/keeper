<template>
  <div class="container">
    <div class="header">
      <FontIcon icon="img macro" />
      <div class="ml-2">
        {{selectedMacro?.title}}
      </div>
    </div>
    <div class="buttons">
      <ToolbarButton icon="icon run">Execute</ToolbarButton>
      <ToolbarButton icon="icon close">Close</ToolbarButton>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, inject, toRef} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ToolbarButton from '/@/second/buttons/ToolbarButton.vue'
export default defineComponent({
  name: "MacroHeader",
  components: {
    FontIcon,
    ToolbarButton
  },
  props: {
    onExecute: {
      type: Function as PropType<() => void>
    }
  },
  setup(props) {
    const selectedMacro = inject('selectedMacro') as any
    const onExecute = toRef(props, 'onExecute')
    return {
      onExecute,
      selectedMacro
    }
  }
})
</script>

<style scoped>
.container {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  background: var(--theme-bg-modalheader);
  height: var(--dim-toolbar-height);
  min-height: var(--dim-toolbar-height);
  overflow: hidden;
  border-top: 1px solid var(--theme-border);
  border-bottom: 1px solid var(--theme-border);
}

.header {
  font-weight: bold;
  margin-left: 10px;
  display: flex;
  align-items: center;
}

.buttons {
  display: flex;
  align-items: stretch;
}
</style>
