<template>
  <!--  <LoadingInfo wrapper message="Waiting for structure"/>-->
  <!--  <ErrorInfo :message="errorMessage" alignTop/>-->
  <div class="container" ref="container" @wheel="handleGridWheel">
    <input/>
    <table
      class="table"
      @mousedown="handleGridMouseDown"
      @mousemove="handleGridMouseMove"
      @mouseup="handleGridMouseUp"
    >
      <thead>
      <tr>
        <td
          class="header-cell"
          data-row="header"
          data-col="header"
          :style="`width:${headerColWidth}px; min-width:${headerColWidth}px; max-width:${headerColWidth}px`"
        >
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
      <tr v-if="display && display.filterable">
        <td
          class="header-cell"
          data-row="filter"
          data-col="header"
          :style="`width:${headerColWidth}px; min-width:${headerColWidth}px; max-width:${headerColWidth}px`">
          <InlineButton v-if="display.filterCount > 0" @click="() => display.clearFilters()" square>
            <FontIcon icon="icon filter-off"/>
          </InlineButton>
        </td>

        <td
          v-for="(col, index) in visibleRealColumns"
          class="filter-cell"
          data-row="filter"
          :data-col="col.colIndex"
          :style="`width:${col.width}px; min-width:${col.width}px; max-width:${col.width}px`"
        ></td>
      </tr>

      </thead>

      <tbody>
      <DataGridRow
        v-for="(rowIndex, i) in rowsIndexs"
        :key="i"
        :rowIndex="rowIndex"
        :grider="grider"
        :conid="conid"
        :database="database"
        :visibleRealColumns="visibleRealColumns"
        :rowHeight="rowHeight"
        :autofillSelectedCells="autofillSelectedCells"
        :isDynamicStructure="isDynamicStructure"
        :selectedCells="filterCellsForRow(selectedCells, rowIndex)"
        :autofillMarkerCell="filterCellForRow(autofillMarkerCell, rowIndex)"
        :focusedColumns="display ? display.focusedColumns : null"
        :inplaceEditorState="inplaceEditorState"
        :currentCellColumn="currentCell && currentCell[0] == rowIndex ? currentCell[1] : null"
        :dispatchInsplaceEditor="dispatchInsplaceEditor"
        :frameSelection="frameSelection"
      />
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, inject, nextTick, PropType, ref, toRefs, unref, watch} from 'vue'
import {compact, isEqual, isNaN, isNumber, max, range, sumBy, uniq} from 'lodash-es'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import CollapseButton from '/@/second/datagrid/CollapseButton.vue'
import ColumnHeaderControl from '/@/second/datagrid/ColumnHeaderControl.vue'
import DataGridRow from '/@/second/datagrid/DataGridRow.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {
  cellIsSelected,
  countColumnSizes,
  countVisibleRealColumns,
  filterCellForRow,
  filterCellsForRow
} from '/@/second/datagrid/gridutil'
import {GridDisplay} from "/@/second/keeper-datalib";
import Grider from "/@/second/datagrid/Grider";
import {SeriesSizes} from "/@/second/datagrid/SeriesSizes";
import {
  CellAddress,
  cellFromEvent,
  emptyCellArray,
  getCellRange,
  isRegularCell,
  nullCell,
  topLeftCell
} from './selection'
import createRef from '/@/second/utility/createRef'
import {isCtrlOrCommandKey} from '/@/second/utility/common'
import createReducer from '/@/second/utility/createReducer'

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
    InlineButton,
    FontIcon,
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
    const {errorMessage, grider, display, onLoadNextData, collapsedLeftColumnStore} = toRefs(props)
    const container = ref<Nullable<HTMLElement>>(null)

    const wheelRowCount = ref(5)
    const tabVisible = inject('tabVisible')
    const containerWidth = ref(503)
    const containerHeight = ref(908)
    const rowHeight = computed(() => 25) //todo  $: rowHeight = $dataGridRowHeight;

    const firstVisibleRowScrollIndex = ref(0)
    const firstVisibleColumnScrollIndex = ref(0)


    const currentCell = ref<CellAddress>(topLeftCell)
    const selectedCells = ref<CellAddress[]>([topLeftCell])
    const dragStartCell = ref<CellAddress | null>(nullCell)
    const shiftDragStartCell = ref<CellAddress | null>(nullCell)
    const autofillDragStartCell = ref<CellAddress | null>(nullCell)
    const autofillSelectedCells = ref<CellAddress[]>(emptyCellArray)
    const domFilterControlsRef = createRef<object>({})
    const tabid = inject('tabid')


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

    const autofillMarkerCell = computed(() => selectedCells.value && selectedCells.value.length > 0 && uniq(selectedCells.value.map(x => x[0])).length == 1
      ? [max(selectedCells.value.map(x => x[0])), max(selectedCells.value.map(x => x[1]))]
      : null)

    const columns = computed(() => display.value?.allColumns || [])

    const columnSizes = ref<SeriesSizes>()
    watch(() => [grider.value, columns.value, containerWidth.value, display.value], async () => {
      await nextTick()
      columnSizes.value = countColumnSizes(grider.value!, columns.value, containerWidth.value, display.value!)
    })
    // const columnSizes = computed(() => _columnSizes)
    const headerColWidth = computed(() => 40)

    const gridScrollAreaHeight = computed(() => containerHeight.value - 2 * rowHeight.value)
    const gridScrollAreaWidth = computed(() => 205)
    // const gridScrollAreaWidth = computed(() => columnSizes.value ? containerWidth.value - columnSizes.value?.frozenSize - headerColWidth.value - 32 : 0)
    const visibleRowCountUpperBound = computed(() => Math.ceil(gridScrollAreaHeight.value / Math.floor(Math.max(1, rowHeight.value))))
    const visibleRowCountLowerBound = computed(() => Math.floor(gridScrollAreaHeight.value / Math.ceil(Math.max(1, rowHeight.value))))

    const visibleRealColumns = computed(() => columnSizes.value ? countVisibleRealColumns(
      columnSizes.value,
      firstVisibleColumnScrollIndex.value,
      gridScrollAreaWidth.value,
      columns.value,
    ) : [])
    // const visibleRealColumns = computed(() => _visibleRealColumns)


    const selectedCellsInfo = computed(() => getSelectedCellsInfo(selectedCells.value, grider.value!, realColumnUniqueNames.value, getSelectedRowData()))

    const realColumnUniqueNames = computed<any[]>(() => range(columnSizes.value!.realCount).map(
      realIndex => (columns.value[columnSizes.value!.realToModel(realIndex)] || {}).uniqueName
    ))

    const maxScrollColumn = computed(() => (columns.value && columnSizes.value) ?
      columnSizes.value?.scrollInView(0, columns.value.length - 1 - columnSizes.value.frozenCount, gridScrollAreaWidth.value) : 0)

    const [inplaceEditorState, dispatchInsplaceEditor] = createReducer((_, action) => {
      switch (action.type) {
        case 'show':
          if (grider.value && !grider.value.editable) return {}
          return {
            cell: action.cell,
            text: action.text,
            selectAll: action.selectAll,
          }
        case 'close':

          if (action.mode == 'enter' || action.mode == 'tab' || action.mode == 'shiftTab') {
            setTimeout(() => {
              if (isRegularCell(currentCell.value)) {
                switch (action.mode) {
                  case 'enter':
                    break
                  case 'tab':
                    break
                  case 'shiftTab':
                    break
                }
              }
            }, 0)
          }
      }
      return {}
    }, {})

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


    function showMultilineCellEditorConditional(cell) {
      if (!cell) return false
      const rowData = grider.value!.getRowData(cell[0])
      if (!rowData) return null
      const cellData = rowData[realColumnUniqueNames.value[cell[1]]]
      //todo
      /*if (shouldOpenMultilineDialog(cellData)) {
        showModal(EditCellDataModal, {
          value: cellData,
          onSave: value => grider.setCellValue(cell[0], realColumnUniqueNames[cell[1]], value),
        });
        return true;
      }*/
      return false
    }

    function handleGridWheel(event) {
      if (event.shiftKey) {
        scrollHorizontal(event.deltaY, event.deltaX);
      } else {
        scrollHorizontal(event.deltaX, event.deltaY);
        scrollVertical(event.deltaX, event.deltaY);
      }
    }

    function scrollVertical(deltaX, deltaY) {

    }

    function scrollHorizontal(deltaX, deltaY) {

    }

    function handleGridMouseDown(event) {
      if (event.target.closest('.buttonLike')) return
      if (event.target.closest('.resizeHandleControl')) return
      if (event.target.closest('.collapseButtonMarker')) return
      if (event.target.closest('.showFormButtonMarker')) return
      if (event.target.closest('input')) return

      shiftDragStartCell.value = null
      event.preventDefault()

      const cell = cellFromEvent(event)
      if (event.button == 2) {
        if (cell && !cellIsSelected(cell[0], cell[1], selectedCells.value)) {
          selectedCells.value = [cell]
        }
        return
      }

      const autofill = event.target.closest('div.autofillHandleMarker')
      if (autofill) {
        autofillDragStartCell.value = cell;
      } else {
        const oldCurrentCell = currentCell.value;
        currentCell.value = cell;

        if (isCtrlOrCommandKey(event)) {
          if (isRegularCell(cell)) {
            if (selectedCells.value.find(x => x[0] == cell[0] && x[1] == cell[1])) {
              selectedCells.value = selectedCells.value.filter(x => x[0] != cell[0] || x[1] != cell[1]);
            } else {
              selectedCells.value = [...selectedCells.value, cell];
            }
          }
        } else if (event.shiftKey) {
          selectedCells.value = getCellRange(oldCurrentCell, cell);
        } else {
          selectedCells.value = getCellRange(cell, cell);
          dragStartCell.value = cell;

          if (isRegularCell(cell) && !isEqual(cell, inplaceEditorState.value.cell) && isEqual(cell, oldCurrentCell)) {
            if (!showMultilineCellEditorConditional(cell)) {
              dispatchInsplaceEditor({type: 'show', cell, selectAll: true});
            }
          } else if (!isEqual(cell, inplaceEditorState.value.cell)) {
            dispatchInsplaceEditor({type: 'close'});
          }
        }
      }
    }

    function handleGridMouseMove() {

    }

    function handleGridMouseUp() {

    }

    return {
      ...toRefs(props),
      errorMessage,
      columnSizes,
      headerColWidth,
      collapsedLeftColumnStore,
      rowHeight,
      currentCell,
      autofillSelectedCells,
      selectedCells,
      autofillMarkerCell,
      handleGridWheel,
      updateCollapsedLeftColumn,
      filterCellsForRow,
      filterCellForRow,
      visibleRowCountUpperBound,
      visibleRowCountLowerBound,
      visibleRealColumns,
      updateResizeSplitter,
      containerWidth,
      containerHeight,
      inplaceEditorState,
      dispatchInsplaceEditor,
      firstVisibleRowScrollIndex,
      handleGridMouseDown,
      handleGridMouseMove,
      handleGridMouseUp,
      rowsIndexs: range(firstVisibleRowScrollIndex.value, Math.min(firstVisibleRowScrollIndex.value + visibleRowCountUpperBound.value, grider.value!.rowCount))
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

