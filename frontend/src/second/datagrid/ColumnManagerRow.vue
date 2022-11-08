<template>
  <div class="row" @click="e => change(e)" :class="isSelected && 'isSelected'">
    <span class="expandColumnIcon">
       <FontIcon/>
    </span>
    <FontIcon v-if="isJsonView" icon="img column"/>
    <a-input v-else type="checkbox" :checked="column && column.isChecked"/>
    <ColumnLabel v-bind="{...column}" showDataType :conid="conid" :database="database"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef, unref} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'

export default defineComponent({
  name: "ColumnManagerRow",
  components: {
    FontIcon,
    ColumnLabel
  },
  props: {
    column: {
      type: Object as PropType<{
        uniqueName: string
        isChecked: boolean
        uniquePath: string
      }>,
    },
    isJsonView: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    conid: {
      type: String as PropType<string>,
    },
    database: {
      type: String as PropType<string>,
    },
  },
  setup(props) {
    const isJsonView = toRef(props, 'isJsonView')

    function change(e: Event) {
      if ((e.target as HTMLElement).closest('.expandColumnIcon')) return
      if (unref(isJsonView)) {

      }

    }

    return {
      ...props,
      change
    }
  }
})
</script>

<style scoped>
.row {
  margin-left: 5px;
  margin-right: 5px;
  cursor: pointer;
  white-space: nowrap;
}

.row:hover {
  background: var(--theme-bg-hover);
}

.row.isSelected {
  background: var(--theme-bg-selected);
}
</style>
