<template>
  <tr :style="`height: ${rowHeight}px`">
    <RowHeaderCell
      :rowIndex="rowIndex"
      :showForm="setFormView ? () => setFormView(rowData, null) : null"/>
    <template v-for="(col, index) in  visibleRealColumns" :key="index">
      <td class="editor"
          v-if="inplaceEditorState
          && inplaceEditorState.cell &&
          rowIndex == inplaceEditorState.cell[0] &&
          col.colIndex == inplaceEditorState.cell[1]">
        column deitor 输入
      </td>
      <DataGridCell
        v-else
        :rowIndex="rowIndex"
        :rowData="rowData"
        :col="col"
        :conid="conid"
        :database="database"
      />
    </template>
  </tr>
</template>

<script lang="ts">
import {isNumber} from 'lodash-es'
import {computed, defineComponent, PropType, toRefs} from 'vue'
import RowHeaderCell from '/@/second/datagrid/RowHeaderCell.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'
import Grider from '/@/second/datagrid/Grider'
import {MacroSelectedCell} from '/@/second/keeper-datalib'
import {CellAddress} from './selection'

export default defineComponent({
  name: "DataGridRow",
  components: {
    RowHeaderCell,
    DataGridCell
  },
  props: {
    rowHeight: {
      type: Number as PropType<number>,
    },
    rowIndex: {
      type: Number as PropType<number>,
    },
    visibleRealColumns: {
      type: Array as PropType<any[]>
    },
    grider: {
      type: Object as PropType<Grider>,
    },
    frameSelection: {
      type: Boolean as PropType<boolean>
    },
    selectedCells: {
      type: Array as PropType<MacroSelectedCell[]>
    },
    autofillMarkerCell: {
      type: Array as PropType<CellAddress[]>
    },
    inplaceEditorState: {
      type: Object as PropType<{ [key in string]: unknown }>,
    },
    setFormView: {
      type: Function as PropType<Function>,
    },
    dispatchInsplaceEditor: {
      type: Function as PropType<(action: any) => void>
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    currentCellColumn: {
      type: [Number, String] as PropType<number | 'header' | 'filter'>,
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    focusedColumns: {
      type: Array as PropType<string[]>
    },
  },
  setup(props) {
    const {grider, rowIndex} = toRefs(props)

    const rowData = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowData(rowIndex.value)
      }
      return null
    })

    return {
      ...toRefs(props),
      rowData,
    }
  }
})
</script>

<style scoped>
tr {
  background-color: var(--theme-bg-0);
}

td.editor {
  position: relative;
}

tr:nth-child(6n + 3) {
  background-color: var(--theme-bg-1);
}

tr:nth-child(6n + 6) {
  background-color: var(--theme-bg-alt);
}
</style>
