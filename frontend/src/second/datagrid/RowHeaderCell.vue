<template>
  <td
    :data-row="`${rowIndex}`"
    data-col="header"
    @mouseenter="() => (mouseIn = true)"
    @mouseleave="() => (mouseIn = false)">
    {{ rowIndex + 1 || '' }}
    <ShowFormButton v-if="mouseIn && showForm" @click.stop.prevent="showForm"/>
  </td>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRefs} from 'vue'
import ShowFormButton from '/@/second/formview/ShowFormButton.vue'

export default defineComponent({
  name: "RowHeaderCell",
  components: {
    ShowFormButton
  },
  props: {
    rowIndex: {
      type: Number as PropType<number>,
    },
    showForm: {
      type: Function as PropType<() => void>,
    }
  },
  setup(props) {
    return {
      ...toRefs(props),
      mouseIn: ref(false),
    }
  }
})
</script>

<style scoped>
td {
  border: 1px solid var(--theme-border);
  text-align: left;
  padding: 0 2px;
  background-color: var(--theme-bg-1);
  overflow: hidden;
  position: relative;
}
</style>
