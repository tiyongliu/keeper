<template>
  <div class="m-1">
    <div class="space-between">
      <ColumnLabel :columnName="uniqueName"/>
      <InlineButton square narrow @click="handler">
        <FontIcon icon="icon close"/>
      </InlineButton>
    </div>
  </div>

</template>

<script lang="ts">
import {defineComponent, PropType, toRefs, unref} from 'vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {getFilterType} from '/@/second/keeper-filterparser'
import {GridDisplay} from '/@/second/keeper-datalib'

export default defineComponent({
  name: 'JsonViewFilterColumn',
  components: {
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
    const {uniqueName, display} = toRefs(props)

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

    return {
      uniqueName,
      display,
      handler,
      computeFilterType
    }
  }
})
</script>
