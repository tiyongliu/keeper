<template>
  <div class="container">
    <div class="header">
      <FontIcon icon="img reference"/>
      <div class="ml-2" v-if="reference">
        {{ reference.pureName }} [{{ reference.columns.map(x => x.refName).join(', ') || '' }}] =
        master [
        {{ reference.columns.map(x => x.baseName).join(', ') }}
      </div>
    </div>
    <ToolbarButton icon="icon close" @click="handleClose">Close</ToolbarButton>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ToolbarButton from '/@/second/buttons/ToolbarButton.vue'

export default defineComponent({
  name: "ReferenceHeader",
  components: {
    FontIcon,
    ToolbarButton
  },
  props: {
    reference: {
      type: Object as PropType<{
        pureName: string
        columns: { refName: string; baseName: string }[]
      }>,
    },
  },
  emits: ['close'],
  setup(props, {emit}) {
    function handleClose() {
      emit('close')
    }

    const reference = toRef(props, 'reference')

    return {
      reference,
      handleClose
    }
  }
})
</script>

<style scoped>
.container {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
}
</style>
