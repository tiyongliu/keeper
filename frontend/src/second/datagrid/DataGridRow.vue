<template>
  <tr :style="`height: ${rowHeight}px`">
<!--
:showForm="setFormView ? () => setFormView(rowData, null) : null"
-->
    <RowHeaderCell :rowIndex="rowIndex" :showForm="handle"/>
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
import {CellAddress} from './selection'
import Grider from "/@/second/datagrid/Grider"
import {MacroSelectedCell} from "/@/second/keeper-datalib";

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
      type: Function as PropType<() => any>,
    },
    dispatchInsplaceEditor: {
      type: Function as PropType<(action: any) => void>
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    currentCellColumn: {
      type: Number as PropType<number>,
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    focusedColumns: {
      type: Array as PropType<string[]>
    }
  },
  emits: ['setFormViewTest'],
  setup(props, context) {
    const {grider, rowIndex} = toRefs(props)
    // const {onShowForm} = attrs

    const rowData = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowData(rowIndex.value)
      }
      return null
    })

    function handle() {
      console.log(`rowData`, context)
    }
    return {
      ...toRefs(props),
      rowData,
      handle
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
