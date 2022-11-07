<template>
  <LoadingInfo
    v-if="!display || (!isDynamicStructure && (!columns || columns.length == 0))"
    wrapper message="Waiting for structure"/>

  <div v-else-if="errorMessage">
    <ErrorInfo :message="errorMessage" alignTop/>
  </div>

  <div v-else-if="isDynamicStructure && isLoadedAll && grider && grider?.rowCount == 0">
    <ErrorInfo
      alignTop
      :message="grider.editable ? 'No rows loaded, check filter or add new documents. You could copy documents from ohter collections/tables with Copy advanved/Copy as JSON command.'
        : 'No rows loaded'"
    />
  </div>

  <div v-else-if="grider && grider.errors && grider.errors.length > 0">
    <ErrorInfo v-for="(err, key) in grider.errors" :key="key" :message="err" isSmall/>
  </div>

  <div v-else class="container" ref="container" @wheel="handleGridWheel">
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
        >
          <DataFilterControl
            :foreignKey="col.foreignKey"
            :columnName="col.uniquePath.length == 1 ? col.uniquePath[0] : null"
            :uniqueName="col.uniqueName"
            :pureName="col.pureName"
            :schemaName="col.schemaName"
            :conid="conid"
            :database="database"
            :jslid="jslid"
            :driver="display?.driver"
            :filterType="useEvalFilters ? 'eval' : col.filterType || getFilterType(col.dataType)"
            :filter="display.getFilter(col.uniqueName)"
            :setFilter="value => display.setFilter(col.uniqueName, value)"
            showResizeSplitter
            :resizeSplitter="(e) => display.resizeColumn(col.uniqueName, col.width, e.detail)"
          />
        </td>
      </tr>
      </thead>

      <tbody>
      <DataGridRow
        v-for="(rowIndex, i) in
        (grider ? range(firstVisibleRowScrollIndex, Math.min(firstVisibleRowScrollIndex + visibleRowCountUpperBound, grider.rowCount)) : [])"
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
    <HorizontalScrollBar
      :minimum="0"
      :maximum="maxScrollColumn"
      :viewportRatio="(gridScrollAreaWidth && columnSizes) ? gridScrollAreaWidth / columnSizes.getVisibleScrollSizeSum() : null"
      @dispatchScroll="e => firstVisibleColumnScrollIndex = e"
      ref="domHorizontalScroll"
    />
    <VerticalScrollBar
      :minimum="0"
      :maximum="grider ? grider.rowCount - visibleRowCountUpperBound + 2 : null"
      :viewportRatio="grider ? visibleRowCountUpperBound / grider.rowCount : null"
      @dispatchScroll="e => firstVisibleRowScrollIndex = e"
      ref="domVerticalScroll"
    />
    <div v-if="selectedCellsInfo" class="row-count-label">{{selectedCellsInfo}}</div>
    <div v-else-if="allRowCount != null && multipleGridsOnTab" class="row-count-label">
      Rows: {allRowCount.toLocaleString()}
    </div>

    <LoadingInfo v-if="isLoading" wrapper message="Loading data" />

  </div>
</template>

<script lang="ts">
import {computed, defineComponent, inject, nextTick, PropType, ref, toRefs, unref, watch} from 'vue'
import {compact, isEqual, isNaN, isNumber, max, range, sumBy, uniq} from 'lodash-es'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import CollapseButton from '/@/second/datagrid/CollapseButton.vue'
import HorizontalScrollBar from '/@/second/datagrid/HorizontalScrollBar.vue'
import VerticalScrollBar from '/@/second/datagrid/VerticalScrollBar.vue'
import ColumnHeaderControl from '/@/second/datagrid/ColumnHeaderControl.vue'
import DataFilterControl from '/@/second/datagrid/DataFilterControl.vue'
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
import {useStatusBarTabItem} from '/@/second/widgets/useStatusBarTabItem'
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
import {getFilterType} from '/@/second/keeper-filterparser'
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
    DataFilterControl,
    HorizontalScrollBar,
    VerticalScrollBar,
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
    },
    jslid: {
      type: [String, Number] as PropType<string | number>
    }
  },
  setup(props) {
    const {errorMessage, grider, display, onLoadNextData, collapsedLeftColumnStore, allRowCount} = toRefs(props)
    //StatusBarTabItem hooks
    useStatusBarTabItem(allRowCount)
    const container = ref<Nullable<HTMLElement>>(null)
    const domHorizontalScroll = ref<Nullable<{scroll: (value: number) => void}>>(null)
    const domVerticalScroll = ref<Nullable<{scroll: (value: number) => void}>>(null)

    const wheelRowCount = ref(5)
    const tabVisible = inject('tabVisible')

    const containerWidth = computed(() => container.value ? container.value.clientWidth : 0)
    const containerHeight = computed(() => container.value ? container.value.clientHeight : 0)
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
    const gridScrollAreaWidth = computed(() => columnSizes.value ? containerWidth.value - columnSizes.value?.frozenSize - headerColWidth.value - 32 : 0)

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
      let newFirstVisibleColumnScrollIndex = firstVisibleColumnScrollIndex.value
      if (deltaX > 0 && deltaY === -0) {
        newFirstVisibleColumnScrollIndex++;
      } else if (deltaX < 0 && deltaY === -0) {
        newFirstVisibleColumnScrollIndex--;
      }

      if (newFirstVisibleColumnScrollIndex > maxScrollColumn.value) {
        newFirstVisibleColumnScrollIndex = maxScrollColumn.value
      }
      if (newFirstVisibleColumnScrollIndex < 0) {
        newFirstVisibleColumnScrollIndex = 0
      }
      firstVisibleColumnScrollIndex.value = newFirstVisibleColumnScrollIndex
      domHorizontalScroll.value!.scroll(newFirstVisibleColumnScrollIndex)
    }
    
    function scrollIntoView(cell) {
      const [row, col] = cell;

      if (row != null) {
        let newRow: number | null = null
        const rowCount = grider.value!.rowCount
        if (rowCount == 0) return

        if (row < firstVisibleRowScrollIndex) newRow = row
        else if (row + 1 >= firstVisibleRowScrollIndex.value + visibleRowCountLowerBound.value)
          newRow = row - visibleRowCountLowerBound.value + 2

        if (newRow! < 0) newRow = 0
        if (newRow! >= rowCount) newRow = rowCount - 1

        if (newRow != null) {
          firstVisibleRowScrollIndex.value = newRow
          domVerticalScroll.value!.scroll(newRow)
        }
      }

      if (col != null) {
        if (col >= columnSizes.value!.frozenCount) {
          let newColumn = columnSizes.value!.scrollInView(
            firstVisibleColumnScrollIndex.value,
            col - columnSizes.value!.frozenCount,
            gridScrollAreaWidth.value
          )
          firstVisibleColumnScrollIndex.value = newColumn;

          domHorizontalScroll.value!.scroll(newColumn);
        }
      }
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
      container,
      domHorizontalScroll,
      domVerticalScroll,
      ...toRefs(props),
      errorMessage,
      selectedCellsInfo,
      columns,
      columnSizes,
      headerColWidth,
      collapsedLeftColumnStore,
      rowHeight,
      currentCell,
      autofillSelectedCells,
      selectedCells,
      autofillMarkerCell,
      getFilterType,
      handleGridWheel,
      updateCollapsedLeftColumn,
      filterCellsForRow,
      filterCellForRow,
      maxScrollColumn,
      gridScrollAreaWidth,
      visibleRowCountUpperBound,
      visibleRowCountLowerBound,
      visibleRealColumns,
      updateResizeSplitter,
      containerWidth,
      containerHeight,
      inplaceEditorState,
      dispatchInsplaceEditor,
      firstVisibleRowScrollIndex,
      firstVisibleColumnScrollIndex,
      handleGridMouseDown,
      handleGridMouseMove,
      handleGridMouseUp,
      range,
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

