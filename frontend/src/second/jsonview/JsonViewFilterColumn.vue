<template>
  <div class="m-1">
    <div class="space-between">
      <ColumnLabel :columnName="uniqueName"/>
      <InlineButton square narrow @click="handler">
        <FontIcon icon="icon close"/>
      </InlineButton>
    </div>
    <DataFilterControl
      :filterType="dynamicFilterType"
      :filter="dynamicFilter"
      :setFilter="handlerSetFilter"
    />
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, toRefs, unref} from 'vue'
import DataFilterControl from '/@/second/datagrid/DataFilterControl.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {getFilterType} from '/@/second/keeper-filterparser'
import {GridDisplay} from '/@/second/keeper-datalib'

export default defineComponent({
  name: 'JsonViewFilterColumn',
  components: {
    DataFilterControl,
    ColumnLabel,
    InlineButton,
    FontIcon,
  },
  props: {
    uniqueName: {
      type: String as PropType<string>,
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    filters: {
      type: Object as PropType<{ [uniqueName: string]: string }>
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    useEvalFilters: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  setup(props) {
    const {uniqueName, display, filters, isDynamicStructure, useEvalFilters} = toRefs(props)
    const dynamicFilter = computed(() => (filters.value && uniqueName.value) ? filters.value[uniqueName.value] : null)
    const dynamicFilterType = computed(() => computeFilterType(isDynamicStructure.value, display.value, uniqueName.value, useEvalFilters.value))

    function handler() {
      (display.value && uniqueName.value) && display.value.removeFilter(uniqueName.value)
    }

    function computeFilterType(isDynamicStructure, display, uniqueName, useEvalFilters) {
      if (unref(useEvalFilters)) return 'eval'
      if (unref(isDynamicStructure)) return 'mongo'
      const col = unref(display) && unref(display).findColumn(uniqueName)
      if (col) {
        return col.filterType || getFilterType(col.dataType)
      }
      return 'string'
    }

    function handlerSetFilter(value: any) {
      display.value && display.value.setFilter(uniqueName.value, unref(value))
    }

    return {
      uniqueName,
      display,
      dynamicFilter,
      dynamicFilterType,
      handler,
      computeFilterType,
      handlerSetFilter
    }
  }
})
</script>
