<template>
  <td :data-row="`${rowIndex}`" data-col="header" @mouseenter="() => (mouseIn = true)" @mouseleave="() => (mouseIn = false)">
    {{rowIndex + 1 || ''}}
    <ShowFormButton v-if="mouseIn && onShowForm" @click="onShowForm"/>
  </td>
</template>

<script lang="ts">
import {defineComponent, PropType, ref} from 'vue'
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
    onShowForm: {
      type: Function as PropType<() => any>,
    }
  },
  setup(props) {
    const mouseIn = ref(false)
    return {
      ...props,
      mouseIn
    }
  }
})
</script>

<style scoped>
td {
  border: 1px solid var(--theme-border);
  text-align: left;
  padding: 2px;
  background-color: var(--theme-bg-1);
  overflow: hidden;
  position: relative;
}
</style>
