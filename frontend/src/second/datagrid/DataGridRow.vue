<template>
  <tr :style="`height: ${rowHeight}px`">
    <RowHeaderCell :rowIndex="rowIndex"/>
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
import {defineComponent, PropType, toRefs, computed, onMounted} from 'vue'
import RowHeaderCell from '/@/second/datagrid/RowHeaderCell.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'
import {CellAddress} from './selection'
import Grider from "/@/second/datagrid/Grider"

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
      type: Function as PropType<() => CellAddress[] | null>
    },
    autofillMarkerCell: {
      type: Function as PropType<() => CellAddress[] | null>
    },
    inplaceEditorState: {
      type: Object as PropType<{ [key in string]: unknown }>,
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
  setup(props) {
    const {grider, rowIndex} = toRefs(props)
    // (grider.value && rowIndex.value) ?  : null
    // const rowData = computed(() => {
    const rowData = computed(() => {
      if (grider.value && rowIndex.value) {
        console.log(`11111111111111111111111111`, rowIndex.value)
        return grider.value.getRowData(rowIndex.value)
      }

      return null
    })

    onMounted(() => {
      setTimeout(() => {
        console.log(`????`, rowData.value)
      }, 5000)
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
