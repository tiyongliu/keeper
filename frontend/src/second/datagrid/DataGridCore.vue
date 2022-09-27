<template>
<!--  <LoadingInfo wrapper message="Waiting for structure"/>-->
<!--  <ErrorInfo :message="errorMessage" alignTop/>-->
  <div
    class="container"
    ref="container"
  >
    <input/>
    <table class="table  sdfssdfsdfsdfsfsd">
      <thead>
      <tr>
        <td class="header-cell" data-row="header" data-col="header"
            :style="`width:${headerColWidth}px; min-width:${headerColWidth}px; max-width:${headerColWidth}px`">
          <CollapseButton :collapsed="collapsedLeftColumnStore" @click="updateCollapsedLeftColumn"/>
        </td>
        <td
          v-for="(col, index) in visibleRealColumns"
          class="header-cell"
          data-row="header"
          :data-col="col.colIndex"
          :key="index"
          :style="`width:${col.width}px; min-width:${col.width}px; max-width:${col.width}px`">
          <ColumnHeaderControl
            :column="col"
            :conid="conid"
            :database="database"
            @resizeSplitter="e => updateResizeSplitter"
            :allowDefineVirtualReferences="allowDefineVirtualReferences"
          />
        </td>
      </tr>
      </thead>

      <tbody>
      <DataGridRow
        v-for="(rowIndex, i) in [0, -1]"
        :key="i"
        :rowIndex="rowIndex"
      />
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs, unref, onMounted, nextTick, } from 'vue'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import CollapseButton from '/@/second/datagrid/CollapseButton.vue'
import ColumnHeaderControl from '/@/second/datagrid/ColumnHeaderControl.vue'
import DataGridRow from '/@/second/datagrid/DataGridRow.vue'
import {countColumnSizes, countVisibleRealColumns} from '/@/second/datagrid/gridutil'
import _columns from './columns.json'
import _columnSizes from './columnSizes.json'
import _visibleRealColumns from './visibleRealColumns.json'
import {TableGridDisplay} from "/@/second/keeper-datalib";
export default defineComponent({
  name: 'DataGridCore',
  components: {
    ErrorInfo,
    LoadingInfo,
    CollapseButton,
    ColumnHeaderControl,
    DataGridRow,
  },
  props: {
    grider: {
      type: Object as PropType<{[key in string]: unknown}>,
      default: undefined
    },
    errorMessage: {
      type: String as PropType<string>,
    },
    collapsedLeftColumnStore: {
      type: Boolean as PropType<boolean>,
      default: true
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    allowDefineVirtualReferences: {
      type: Boolean as PropType<boolean>,
      default: false
    },


    display: {
      type: Object as PropType<TableGridDisplay>
    },
  },
  setup(props) {
    const container = ref<Nullable<HTMLElement>>(null)
    const firstVisibleColumnScrollIndex = ref(0)
    const containerWidth = ref(310)
    const containerHeight = ref(618)
    const containerWidth = ref(0)
    const {errorMessage, display} = toRefs(props)

    const columnSizes = computed(() => _columnSizes)


    const {collapsedLeftColumnStore} = toRefs(props)
    const headerColWidth = computed(() => 40)

    const gridScrollAreaWidth = computed(() => 205)
    const columns = computed(() => _columns)
    // const visibleRealColumns = computed(() => countVisibleRealColumns(
    //   unref(columnSizes),
    //   unref(firstVisibleColumnScrollIndex),
    //   unref(gridScrollAreaWidth),
    //   unref(columns),
    // ))

    const visibleRealColumns = computed(() => _visibleRealColumns)

    function updateCollapsedLeftColumn() {
      collapsedLeftColumnStore.value = !unref(collapsedLeftColumnStore)
    }

    function updateResizeSplitter() {

    }

    const columns = computed(() => display.value?.allColumns || [])
    // countColumnSizes()
    onMounted(async () => {
      await nextTick()

      // console.log(visibleRealColumns, `visibleRealColumnsvisibleRealColumnsvisibleRealColumnsvisibleRealColumnsvisibleRealColumns`)

      // container.value!.clientWidth = 310
      // container.value!.clientHeight = 680
      // container.value!.clientWidth = clientWidth.value
      // container.value!.clientHeight = clientHeight.value
    })
    return {
      ...toRefs(props),
      errorMessage,
      columnSizes,
      headerColWidth,
      collapsedLeftColumnStore,
      updateCollapsedLeftColumn,
      visibleRealColumns,
      updateResizeSplitter,
      containerWidth,
      containerHeight,
    }
  }
})
</script>

<style scoped>
.container {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  user-select: none;
  overflow: hidden;
}

.table {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 20px;
  overflow: scroll;
  border-collapse: collapse;
  outline: none;
}

.header-cell {
  border: 1px solid var(--theme-border);
  text-align: left;
  padding: 0;
  margin: 0;
  background-color: var(--theme-bg-1);
  overflow: hidden;
}

.filter-cell {
  text-align: left;
  overflow: hidden;
  margin: 0;
  padding: 0;
}

.focus-field {
  position: absolute;
  left: -1000px;
  top: -1000px;
}

.row-count-label {
  position: absolute;
  background-color: var(--theme-bg-2);
  right: 40px;
  bottom: 20px;
}
</style>

