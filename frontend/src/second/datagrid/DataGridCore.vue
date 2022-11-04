<template>
  <!--  <LoadingInfo wrapper message="Waiting for structure"/>-->
  <!--  <ErrorInfo :message="errorMessage" alignTop/>-->
  <div class="container" ref="container">
    <input />
    <table class="table">
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
import {
  computed,
  defineComponent,
  nextTick,
  onMounted,
  PropType,
  ref,
  toRefs,
  unref,
  watch
} from 'vue'
import {isNumber, sumBy, isNaN, range, compact, uniq} from 'lodash-es'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import CollapseButton from '/@/second/datagrid/CollapseButton.vue'
import ColumnHeaderControl from '/@/second/datagrid/ColumnHeaderControl.vue'
import DataGridRow from '/@/second/datagrid/DataGridRow.vue'
import {countColumnSizes, countVisibleRealColumns} from '/@/second/datagrid/gridutil'
import _visibleRealColumns from './visibleRealColumns.json'
import {GridDisplay} from "/@/second/keeper-datalib";
import Grider from "/@/second/datagrid/Grider";
import {SeriesSizes} from "/@/second/datagrid/SeriesSizes";
import {topLeftCell} from './selection'
function getSelectedCellsInfo(selectedCells, grider: Grider, realColumnUniqueNames, selectedRowData) {
  if (selectedCells.length > 1 && selectedCells.every(x => isNumber(x[0]) && isNumber(x[1]))) {
    let sum = sumBy(selectedCells, (cell: string[]) => {
      const row = grider.getRowData(cell[0]);
      if (row) {
        const colName = realColumnUniqueNames[cell[1]];
        if (colName) {
          const data = row[colName];
          if (!data) return 0;
          let num = +data;
          if (isNaN(num)) return 0;
          return num;
        }
      }
      return 0;
    });
    let count = selectedCells.length;
    let rowCount = selectedRowData.length;
    return `Rows: ${rowCount.toLocaleString()}, Count: ${count.toLocaleString()}, Sum:${sum.toLocaleString()}`;
  }
  return null;
}
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
      type: Object as PropType<Grider>,
    },
    display: {
      type: Object as PropType<GridDisplay>,
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    frameSelection: {
      type: Boolean as PropType<boolean>
    },
    isLoading: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    allRowCount: {
      type: Number as PropType<number>,
    },
    focusOnVisible: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    formViewAvailable: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    errorMessage: {
      type: String as PropType<string | null>,
    },
    collapsedLeftColumnStore: {
      type: Boolean as PropType<boolean>,
      default: true
    },
    allowDefineVirtualReferences: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    onLoadNextData: {
      type: Function as PropType<() => void>
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    multipleGridsOnTab: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    tabControlHiddenTab: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    useEvalFilters: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isLoadedAll: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    loadedTime: {
      type: Number as PropType<number>,
      default: 0
    }
  },
  setup(props) {


    const container = ref<Nullable<HTMLElement>>(null)
    const firstVisibleRowScrollIndex = ref(0)
    const firstVisibleColumnScrollIndex = ref(0)
    const containerWidth = ref(503)
    const containerHeight = ref(908)
    const {errorMessage, grider, display, onLoadNextData} = toRefs(props)
    const selectedCells = ref([topLeftCell])

    const {collapsedLeftColumnStore} = toRefs(props)
    const headerColWidth = computed(() => 40)

    const gridScrollAreaWidth = computed(() => 205)
    // const columns = computed(() => _columns)

    function getSelectedRowIndexes() {
      if (selectedCells.value.find(x => x[0] == 'header')) return range(0, grider.value!.rowCount);
      return uniq((selectedCells.value || []).map(x => x[0])).filter(x => isNumber(x));
    }

    function getSelectedColumnIndexes() {
      if (selectedCells.value.find(x => x[1] == 'header')) return range(0, realColumnUniqueNames.value.length);
      return uniq((selectedCells.value || []).map(x => x[1])).filter(x => isNumber(x));
    }

    function getSelectedRowData() {
      return compact(getSelectedRowIndexes().map(index => grider.value!.getRowData(index)));
    }

    function getSelectedColumns() {
      return compact(
        getSelectedColumnIndexes().map((index: number) => ({
          columnName: realColumnUniqueNames.value[index],
        }))
      );
    }

    const columns = computed(() => display.value?.allColumns || [])

    watch(() => [grider.value, columns.value, containerWidth.value, display.value], async () => {
      await nextTick()
      columnSizes.value = countColumnSizes(grider.value!, columns.value, containerWidth.value, display.value!)
      console.log(`ccccccccccccc`, containerWidth.value)
    })


    const columnSizes = ref<SeriesSizes>()
    // const columnSizes = computed(() => _columnSizes)

    // const visibleRealColumns = computed(() => countVisibleRealColumns(
    //   columnSizes.value,
    //   firstVisibleColumnScrollIndex.value,
    //   gridScrollAreaWidth.value,
    //   columns.value,
    // ))

    const visibleRealColumns = computed(() => _visibleRealColumns)



    const selectedCellsInfo = computed(() => getSelectedCellsInfo(selectedCells.value, grider.value!, realColumnUniqueNames.value, getSelectedRowData()))

    const realColumnUniqueNames = computed<any[]>(() => range(columnSizes.value!.realCount).map(
      realIndex => (columns.value[columnSizes.value!.realToModel(realIndex)] || {}).uniqueName
    ))

    const maxScrollColumn = computed(() => columnSizes.value.scrollInView(0, columns.value.length - 1 - columnSizes.value.frozenCount, gridScrollAreaWidth.value))


    function updateCollapsedLeftColumn() {
      collapsedLeftColumnStore.value = !unref(collapsedLeftColumnStore)
    }

    function updateResizeSplitter() {

    }

    // const columns = computed(() => display.value?.allColumns || [])
    // countColumnSizes()



    watch(() => [onLoadNextData.value, display.value], () => {
      if (onLoadNextData.value && display.value) {
        onLoadNextData.value()
      }
    })

    onMounted(() => {
      setTimeout(() => {
        console.log(columnSizes.value, `----------------------------`)
      }, 8888)
    })

    return {
      ...props,
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

