<template>
  <div class="row"
       @click="handleClick"
       @mousedown="handleMousedown"
       @mousemove="handleMousemove"
       @mouseup="handleMouseup"
       :class="isSelected && 'isSelected'">
    <span class="expandColumnIcon"
          :style="column && column.uniquePath && `margin-right: ${5 + (column.uniquePath.length - 1) * 10}px`">
             <FontIcon
               :icon="(column && column.isExpandable && display) ? plusExpandIcon(display.isExpandedColumn(column.uniqueName)) : 'icon invisible-box'"
               @click="handleDisplay"/>
    </span>
    <FontIcon v-if="isJsonView" icon="img column"/>
    <!--    <input
          @click="(e) => {e.stopPropagation()}"
          @mousedown="(e) => {e.preventDefault();e.stopPropagation()}"
          @change="handlerChange"
          v-else type="checkbox" :checked="isChecked"/>-->
    <span
      v-else
      style="margin: 0 3px"
      @click="(e) => {e.stopPropagation()}"
      @mousedown="(e) => {e.stopPropagation()}"
      @change="handlerChange">
      <a-checkbox v-model:checked="isChecked"/>
    </span>
    <ColumnLabel v-bind="{...column}" showDataType :conid="conid" :database="database"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import {isBoolean} from 'lodash-es'
import {plusExpandIcon} from '/@/second/icons/expandIcons'
import {GridDisplay} from '/@/second/keeper-datalib'

import {Checkbox} from 'ant-design-vue'

export default defineComponent({
  name: "ColumnManagerRow",
  components: {
    FontIcon,
    ColumnLabel,
    [Checkbox.name]: Checkbox,
  },
  props: {
    column: {
      type: Object as PropType<{
        uniqueName: string
        isChecked: boolean
        uniquePath: string[]
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
    display: {
      type: Object as PropType<GridDisplay>
    },
  },
  emits: [
    'setVisibility',
    'dispatchMousedown',
    'dispatchMousemove',
    'dispatchMouseup',
    'dispatchClick'],
  setup(props, {emit}) {
    const {column, isJsonView, isSelected, conid, database, display} = toRefs(props)

    const isChecked = ref(false)
    watch(() => column.value, () => {
      if (isBoolean(unref(column)?.isChecked)) {
        isChecked.value = unref(column)?.isChecked!
      }
    }, {immediate: true})

    function handleClick(e: Event) {
      if ((e.target as HTMLElement).closest('.expandColumnIcon')) return
      if (unref(isJsonView)) {
        display.value!.showFilter(column.value!.uniqueName)
      } else {
        display.value!.focusColumns([column.value!.uniqueName]);
      }
      emit('dispatchClick')
    }

    function handlerChange() {
      const newValue = !unref(column)?.isChecked
      display.value!.setColumnVisibility(column.value!.uniquePath, newValue)
      emit('setVisibility', newValue)
    }

    function handleDisplay() {
      display.value && display.value.toggleExpandedColumn(column.value!.uniqueName)
    }

    function handleMousedown(e) {
      emit('dispatchMousedown', e)
    }

    function handleMousemove(e) {
      emit('dispatchMousemove', e)
    }

    function handleMouseup(e) {
      emit('dispatchMouseup', e)
    }

    return {
      column,
      isJsonView,
      isSelected,
      isChecked,
      conid,
      database,
      plusExpandIcon,
      handleClick,
      handlerChange,
      handleDisplay,
      handleMousedown,
      handleMousemove,
      handleMouseup,
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

/*input[type="checkbox"] {
  cursor: default;
  appearance: auto;
  box-sizing: border-box;
  margin: 3px 3px 3px 4px;
  padding: initial;
  vertical-align: middle;
}*/
</style>
