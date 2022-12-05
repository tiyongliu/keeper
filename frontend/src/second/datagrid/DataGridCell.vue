<template>
  <td
    :class="[
      isSelected && 'isSelected',
      isFrameSelected && 'isFrameSelected',
      isModifiedRow && 'isModifiedRow',
      isModifiedCell && 'isModifiedCell',
      isInserted && 'isInserted',
      isDeleted && 'isDeleted',
      isAutofillSelected && 'isAutofillSelected',
      isFocusedColumn && 'isFocusedColumn'
    ]"
    :style="style"
    :data-row="`${rowIndex}`"
    :data-col="`${colIndex == null ? col.colIndex : colIndex}`">
    <CellValue :rowData="rowData" :value="value" :jsonParsedValue="jsonParsedValue"/>
  </td>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, toRefs, unref} from 'vue'
import {get} from 'lodash-es'
import ShowFormButton from '/@/second/formview/ShowFormButton.vue'
import CellValue from '/@/second/datagrid/CellValue.vue'
import {isJsonLikeLongString, safeJsonParse} from '/@/second/keeper-tools'

export default defineComponent({
  name: "DataGridCell",
  components: {
    ShowFormButton,
    CellValue,
  },
  props: {
    rowIndex: {
      type: Number as PropType<number>,
    },
    col: {
      type: Object as PropType<{ colIndex?: number, isStructured: boolean, uniquePath: string, uniqueName: string }>
    },
    rowData: {
      type: Object as PropType<object>
    },
    colIndex: {
      type: Number as PropType<number>,
    },
    allowHintField: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    maxWidth: {
      type: Number as PropType<number>,
    },
    minWidth: {
      type: Number as PropType<number>,
    },
    isSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isFrameSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isModifiedRow: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isModifiedCell: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isInserted: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isDeleted: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isAutofillSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isFocusedColumn: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    //domCell
    showSlot: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    setFormView: {
      type: Function as PropType<Function>,
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isAutoFillMarker: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isCurrentCell: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    // onDictionaryLookup
    // onSetValue
  },
  setup(props) {
    const {col, rowData, maxWidth, minWidth} = toRefs(props)
    const value = computed(() => {
      return unref(col)!.isStructured
        ? get(unref(rowData) || {}, unref(col)!.uniquePath)
        : (unref(rowData) || {})[unref(col)!.uniqueName]
    })

    function computeStyle(col) {
      let res = ''
      if (col.width != null) {
        res += `width:${col.width}px; min-width:${col.width}px; max-width:${col.width}px;`
      } else {
        if (maxWidth.value != null) res += `max-width:${maxWidth.value}px;`
        if (minWidth.value != null) res += `min-width:${minWidth.value}px;`
      }
      return res
    }

    const style = computed(() => computeStyle(col.value))
    const jsonParsedValue = computed(() => isJsonLikeLongString(value.value) ? safeJsonParse(value.value) : null)

    return {
      ...toRefs(props),
      value,
      style,
      jsonParsedValue,
    }
  }
})
</script>

<style scoped>
td {
  font-weight: normal;
  border: 1px solid var(--theme-border);
  padding: 2px;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

td.isFrameSelected {
  outline: 3px solid var(--theme-bg-selected);
  outline-offset: -3px;
}

td.isAutofillSelected {
  outline: 3px solid var(--theme-bg-selected);
  outline-offset: -3px;
}

td.isFocusedColumn {
  background: var(--theme-bg-alt);
}

td.isModifiedRow {
  background: var(--theme-bg-gold);
}

td.isModifiedCell {
  background: var(--theme-bg-orange);
}

td.isInserted {
  background: var(--theme-bg-green);
}

td.isDeleted {
  background: var(--theme-bg-volcano);
}

td.isSelected {
  background: var(--theme-bg-selected);
}

td.isDeleted {
  background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAEElEQVQImWNgIAX8x4KJBAD+agT8INXz9wAAAABJRU5ErkJggg==');
  background-repeat: repeat-x;
  background-position: 50% 50%;
}

.hint {
  color: var(--theme-font-3);
  margin-left: 5px;
}

.autoFillMarker {
  width: 8px;
  height: 8px;
  background: var(--theme-bg-selected-point);
  position: absolute;
  right: 0px;
  bottom: 0px;
  overflow: visible;
  cursor: crosshair;
}
</style>
