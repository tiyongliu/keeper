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
import {defineComponent, PropType, toRefs, computed} from 'vue'
import RowHeaderCell from '/@/second/datagrid/RowHeaderCell.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'
import Grider from "/@/second/datagrid/Grider";

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
    inplaceEditorState: {
      type: Object as PropType<{ [key in string]: unknown }>,
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

    return {
      ...toRefs(props),
      rowData: computed(() => (grider.value && rowIndex.value) ? grider.value.getRowData(rowIndex.value) : null),
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
