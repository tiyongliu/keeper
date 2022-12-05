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
        <InplaceEditor
          :width="col.width"
          :inplaceEditorState="inplaceEditorState"
          :dispatchInsplaceEditor="dispatchInsplaceEditor"
          :cellValue="rowData[col.uniqueName]"
          @setValue="value => grider.setCellValue(rowIndex, col.uniqueName, value)"/>
        <!--        <FormOutlined />
                <CheckOutlined/>
                <CloseOutlined/>-->
      </td>
      <DataGridCell
        v-else
        :rowIndex="rowIndex"
        :rowData="rowData"
        :col="col"
        :conid="conid"
        :database="database"
        :allowHintField="hintFieldsAllowed?.includes(col.uniqueName)"
        :isSelected="frameSelection ? false : cellIsSelected(rowIndex, col.colIndex, selectedCells)"
        :isCurrentCell="col.colIndex == currentCellColumn"
        :isFrameSelected="frameSelection ? cellIsSelected(rowIndex, col.colIndex, selectedCells) : false"
        :isAutofillSelected="cellIsSelected(rowIndex, col.colIndex, autofillSelectedCells)"
        :isFocusedColumn="focusedColumns?.includes(col.uniqueName)"
        :isModifiedCell="rowStatus.modifiedFields && rowStatus.modifiedFields.has(col.uniqueName)"
        :isModifiedRow="rowStatus.status == 'updated'"
        :isInserted="rowStatus.status == 'inserted' ||
          (rowStatus.insertedFields && rowStatus.insertedFields.has(col.uniqueName))"
        :isDeleted="rowStatus.status == 'deleted' ||
          (rowStatus.deletedFields && rowStatus.deletedFields.has(col.uniqueName))"
        :setFormView="setFormView"
        :isDynamicStructure="isDynamicStructure"
        :isAutoFillMarker="autofillMarkerCell &&
          autofillMarkerCell[1] == col.colIndex &&
          autofillMarkerCell[0] == rowIndex &&
          grider.editable"/>
    </template>
  </tr>
</template>

<script lang="ts">
import {isNumber} from 'lodash-es'
import {computed, defineComponent, PropType, toRefs} from 'vue'
import RowHeaderCell from '/@/second/datagrid/RowHeaderCell.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'
import Grider from '/@/second/datagrid/Grider'
import InplaceEditor from './InplaceEditor.vue'
import {MacroSelectedCell} from '/@/second/keeper-datalib'
import {CellAddress} from './selection'
import {cellIsSelected} from './gridutil'

import {Spin} from 'ant-design-vue'
import {CheckOutlined, CloseOutlined, FormOutlined} from '@ant-design/icons-vue';

export default defineComponent({
  name: "DataGridRow",
  components: {
    InplaceEditor,
    Spin,
    RowHeaderCell,
    DataGridCell,

    FormOutlined, CloseOutlined, CheckOutlined,
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
    autofillSelectedCells: {
      type: Array as PropType<CellAddress[]>
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
      type: [String, Number] as PropType<number | 'header' | 'filter'>,
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
    const {grider, rowIndex, visibleRealColumns} = toRefs(props)

    const rowData = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowData(rowIndex.value)
      }
      return null
    })

    const rowStatus = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowStatus(rowIndex.value)
      }
      return null
    })

    const hintFieldsAllowed = computed(() => {
      visibleRealColumns.value ? visibleRealColumns.value?.filter(col => {
        if (!col.hintColumnNames) return false
        if (rowStatus.value && rowStatus.value.modifiedFields && rowStatus.value.modifiedFields.has(col.uniqueName)) return false
        return true;
      }).map(col => col.uniqueName) : []
    })

    return {
      ...toRefs(props),
      rowData,
      rowStatus,
      hintFieldsAllowed,
      cellIsSelected
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
