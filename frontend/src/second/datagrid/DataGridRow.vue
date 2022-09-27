<template>
  <tr :style="`height: ${rowHeight}px`">
    <RowHeaderCell :rowIndex="rowIndex"/>
    <template v-for="(col, index) in  visibleRealColumns" :key="index">
      <DataGridCell
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
import {defineComponent, PropType, toRefs} from 'vue'
import RowHeaderCell from '/@/second/datagrid/RowHeaderCell.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'

export default defineComponent({
  name: "DataGridRow",
  components: {
    RowHeaderCell,
    DataGridCell
  },
  props: {
    rowData: {
      type: Object as PropType<{ [key in string]: unknown }>,
    },
    grider: {
      type: Object as PropType<{ [key in string]: unknown }>,
    },
    rowHeight: {
      type: Number as PropType<number>,
    },
    rowIndex: {
      type: Number as PropType<number>,
    },
    visibleRealColumns: {
      type: Array as PropType<any[]>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
  },
  setup(props) {
    return {
      ...toRefs(props)
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
