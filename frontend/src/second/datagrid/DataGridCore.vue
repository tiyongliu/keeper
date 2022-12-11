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
      :message="grider.editable
      ? 'No rows loaded, check filter or add new documents. You could copy documents from ohter collections/tables with Copy advanved/Copy as JSON command.'
      : 'No rows loaded'"
    />
  </div>

  <div v-else-if="grider && grider.errors && grider.errors.length > 0">
    <ErrorInfo v-for="(err, key) in grider.errors" :key="key" :message="err" isSmall/>
  </div>

  <div v-else class="container" ref="container" @wheel="handleGridWheel">
    <!--  todo 现在还不清楚具体使用场景，暂时注释，如果我们点页面Hide按钮，把列表的所有column隐藏了，就会显示出来  -->
    <input ref="domFocusField" v-if="false"/>
    <table
      class="table"
      @mousedown="handleGridMouseDown"
      @mousemove="handleGridMouseMove"
      @mouseup="handleGridMouseUp">
      <thead>
      <tr>
        <td
          class="header-cell"
          data-row="header"
          data-col="header"
          :style="`width:${headerColWidth}px; min-width:${headerColWidth}px; max-width:${headerColWidth}px`">
          <CollapseButton
            :collapsed="collapsedLeftColumnStore"
            @click="updateCollapsedLeftColumn"/>
        </td>
        <td
          v-for="(col, index) in visibleRealColumns"
          :key="index"
          class="header-cell"
          data-row="header"
          :data-col="`${col.colIndex}`"
          :style="`width:${col.width}px; min-width:${col.width}px; max-width:${col.width}px`">
          <ColumnHeaderControl
            :column="col"
            :conid="conid"
            :database="database"
            :setSort="display && display.sortable ? order => display.setSort(col.uniqueName, order) : null"
            :addToSort="display && display.sortable ? order => display.addToSort(col.uniqueName, order) : null"
            :order="display && display.sortable ? display.getSortOrder(col.uniqueName) : null"
            :orderIndex="display && display.sortable ? display.getSortOrderIndex(col.uniqueName) : -1"
            :isSortDefined="display && display.sortable ? display.isSortDefined() : false"
            :clearSort="display && display.sortable ? () => display.clearSort() : null"
            @resizeSplitter="e => {(display && col) && display.resizeColumn(col.uniqueName, col.width, e.detail)}"
            :setGrouping="display.groupable ? groupFunc => display.setGrouping(col.uniqueName, groupFunc) : null"
            :grouping="display.getGrouping(col.uniqueName)"
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
          :key="index"
          class="filter-cell"
          data-row="filter"
          :data-col="`${col.colIndex}`"
          :style="`width:${col.width}px; min-width:${col.width}px; max-width:${col.width}px`">
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
            :filter="display ? display.getFilter(col.uniqueName) : null"
            :setFilter="value => display.setFilter(col.uniqueName, value)"
            showResizeSplitter
            @dispatchResizeSplitter="(e) => display.resizeColumn(col.uniqueName, col.width, e.detail)"
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
        :focusedColumns="display.focusedColumns"
        :inplaceEditorState="inplaceEditorState"
        :currentCellColumn="currentCell && currentCell[0] == rowIndex ? currentCell[1] : null"
        :dispatchInsplaceEditor="dispatchInsplaceEditor"
        :frameSelection="frameSelection"
        :setFormView="formViewAvailable && display && display?.baseTable?.primaryKey ? handleSetFormView : null"
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
    <div v-if="selectedCellsInfo" class="row-count-label">{{ selectedCellsInfo }}</div>
    <div v-else-if="allRowCount != null && multipleGridsOnTab" class="row-count-label">
      Rows: {allRowCount.toLocaleString()}
    </div>
    <LoadingInfo v-if="isLoading" wrapper message="Loading data"/>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  nextTick,
  PropType,
  ref,
  Ref,
  toRefs,
  unref,
  watch,
  watchEffect
} from 'vue'
import {
  compact,
  flatten,
  isEqual,
  isNaN,
  isNumber,
  max,
  min,
  pick,
  range,
  sumBy,
  uniq
} from 'lodash-es'
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
// import {dataGridRowHeight} from './DataGridRowHeightMeter.vue'
import {GridDisplay} from '/@/second/keeper-datalib'
import Grider from '/@/second/datagrid/Grider'
import {SeriesSizes} from '/@/second/datagrid/SeriesSizes'
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
import openReferenceForm, {openPrimaryKeyForm} from '/@/second/formview/openReferenceForm'
import createReducer from '/@/second/utility/createReducer'
import stableStringify from 'json-stable-stringify'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import keycodes from '/@/second/utility/keycodes'
import bus from '/@/second/utility/bus'

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
    changeSelectedColumns: {
      type: Function as PropType<(cols: any[]) => void>,
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
      type: Object as PropType<Ref<boolean>>,
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
  emits: ['selectedCellsPublished'],
  setup(props, {emit}) {
    const {
      conid,
      database,
      errorMessage,
      grider,
      display,
      onLoadNextData,
      collapsedLeftColumnStore,
      allRowCount,
      changeSelectedColumns,
      focusOnVisible,
    } = toRefs(props)

    const bootstrap = useBootstrapStore()

    //StatusBarTabItem hooks
    useStatusBarTabItem(allRowCount)

    const container = ref<Nullable<HTMLElement>>(null)
    const domFocusField = ref<Nullable<HTMLElement>>(null)
    const domHorizontalScroll = ref<Nullable<{ scroll: (value: number) => void }>>(null)
    const domVerticalScroll = ref<Nullable<{ scroll: (value: number) => void }>>(null)
    const wheelRowCount = ref(5)
    const firstVisibleRowScrollIndex = ref(0)
    const firstVisibleColumnScrollIndex = ref(0)
    const currentCell = ref<CellAddress>(topLeftCell)
    const selectedCells = ref<CellAddress[]>([topLeftCell])
    const dragStartCell = ref<Nullable<CellAddress>>(nullCell)
    const shiftDragStartCell = ref<Nullable<CellAddress>>(nullCell)
    const autofillDragStartCell = ref<Nullable<CellAddress>>(nullCell)
    const autofillSelectedCells = ref<CellAddress[]>(emptyCellArray)
    const columnSizes = ref<SeriesSizes>()

    const domFilterControlsRef = createRef<object>({})
    const lastPublishledSelectedCellsRef = createRef<string>('')

    const containerWidth = ref(container.value ? container.value.clientWidth : 0)
    const containerHeight = ref(container.value ? container.value.clientWidth : 0)

    const columns = computed(() => display.value?.allColumns || [])
    const rowHeight = computed(() => 25) //todo  $: rowHeight = $dataGridRowHeight;
    const autofillMarkerCell = computed(() => selectedCells.value && selectedCells.value.length > 0 && uniq(selectedCells.value.map(x => x[0])).length == 1
      ? [max(selectedCells.value.map(x => x[0])), max(selectedCells.value.map(x => x[1]))]
      : null)
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
    const selectedCellsInfo = computed(() => getSelectedCellsInfo(selectedCells.value, grider.value!, realColumnUniqueNames.value, getSelectedRowData()))
    const realColumnUniqueNames = computed<any[]>(() => (columnSizes.value ? range(columnSizes.value.realCount) : []).map(
      realIndex => (columns.value[columnSizes.value!.realToModel(realIndex)] || {}).uniqueName
    ))
    const maxScrollColumn = computed(() => (columns.value && columnSizes.value) ?
      columnSizes.value?.scrollInView(0, columns.value.length - 1 - columnSizes.value.frozenCount, gridScrollAreaWidth.value) : 0)

    const tabVisible = inject<Ref<Nullable<boolean>>>('tabVisible')
    // const tabid = inject('tabid')

    bus.emitter.on(bus.resize, updateWidgetStyle)

    function updateWidgetStyle() {
      nextTick(() => {
        if (container.value && container.value!.clientWidth) containerWidth.value = container.value!.clientWidth
        if (container.value && container.value!.clientHeight) containerHeight.value = container.value!.clientHeight
      })
    }

    watch(() => [collapsedLeftColumnStore.value], updateWidgetStyle)

    watch(() => [
      // grider.value, columns.value,
      containerWidth.value, display.value], () => {
      columnSizes.value = countColumnSizes(grider.value!, columns.value, containerWidth.value, display.value!)
      updateWidgetStyle()
    })

    // watch(() => [onLoadNextData.value, display.value], () => {
    //   updateWidgetStyle()
    // })

    watch(() => [firstVisibleRowScrollIndex.value, visibleRowCountUpperBound.value], async () => {
      if (onLoadNextData.value && grider.value && firstVisibleRowScrollIndex.value + visibleRowCountUpperBound.value >= grider.value!.rowCount && rowHeight.value > 0) {
        onLoadNextData.value()
      }
    })

    watch(() => selectedCells.value, () => {
      const stringified = stableStringify(selectedCells)
      if (lastPublishledSelectedCellsRef.get() != stringified) {
        lastPublishledSelectedCellsRef.set(stringified)
        const cellsValue = () => getCellsPublished(selectedCells.value)
        emit('selectedCellsPublished', cellsValue)
        bootstrap.subscribeSelectedCellsCallback(cellsValue)
        if (changeSelectedColumns.value) changeSelectedColumns.value(getSelectedColumns().map(x => x.columnName))
      }
    })

    watchEffect(() => {
      if (unref(tabVisible) && domFocusField.value && focusOnVisible.value) {
        domFocusField.value && domFocusField.value.focus()
      }
    })

    function getSelectedRowIndexes() {
      if (selectedCells.value.find(x => x[0] == 'header')) return range(0, grider.value!.rowCount);
      return uniq((selectedCells.value || []).map(x => x[0])).filter(x => isNumber(x));
    }

    function getSelectedColumnIndexes() {
      if (selectedCells.value.find(x => x[1] == 'header')) return range(0, realColumnUniqueNames.value.length);
      return uniq((selectedCells.value || []).map(x => x[1])).filter(x => isNumber(x));
    }

    function getSelectedRowData() {
      return grider.value ? compact(getSelectedRowIndexes().map(index => grider.value!.getRowData(index))) : []
    }

    function getSelectedColumns() {
      return realColumnUniqueNames.value ? compact(
        getSelectedColumnIndexes().map((index: number) => ({
          columnName: realColumnUniqueNames.value[index],
        }))
      ) : []
    }

    const [inplaceEditorState, dispatchInsplaceEditor] = createReducer((_, action) => {
      switch (action.type) {
        case 'show':
          if (!grider.value || !grider.value.editable) return {}
          return {
            cell: action.cell,
            text: action.text,
            selectAll: action.selectAll,
          }
        case 'close':
          if (domFocusField.value) domFocusField.value.focus();
          if (action.mode == 'enter' || action.mode == 'tab' || action.mode == 'shiftTab') {
            setTimeout(() => {
              if (isRegularCell(currentCell.value)) {
                switch (action.mode) {
                  case 'enter':
                    moveCurrentCell(currentCell[0] + 1, currentCell[1]);
                    break
                  case 'tab':
                    moveCurrentCellWithTabKey(false)
                    break
                  case 'shiftTab':
                    moveCurrentCellWithTabKey(true)
                    break
                }
              }
            }, 0)
          }
      }
      return {}
    }, {})

    function handleGridKeyDown(event) {
      if (inplaceEditorState.value) return
      if (
        !event.ctrlKey &&
        !event.altKey &&
        !event.metaKey &&
        ((event.keyCode >= keycodes.a && event.keyCode <= keycodes.z) ||
          (event.keyCode >= keycodes.n0 && event.keyCode <= keycodes.n9) ||
          (event.keyCode >= keycodes.numPad0 && event.keyCode <= keycodes.numPad9) ||
          event.keyCode == keycodes.dash)
      ) {
        // @ts-ignore
        event.preventDefault();
        dispatchInsplaceEditor({type: 'show', text: event.key, cell: currentCell});
      }

      if (event.keyCode == keycodes.f2 || event.keyCode == keycodes.enter) {
        // @ts-ignore
        if (!showMultilineCellEditorConditional(currentCell)) {
          dispatchInsplaceEditor({type: 'show', cell: currentCell, selectAll: true});
        }
      }

      if (event.shiftKey) {
        if (!isRegularCell(shiftDragStartCell.value!)) {
          shiftDragStartCell.value = currentCell.value;
        }
      } else {
        shiftDragStartCell.value = nullCell;
      }

      handleCursorMove(event);

      if (
        event.shiftKey &&
        event.keyCode != keycodes.shift &&
        event.keyCode != keycodes.tab &&
        event.keyCode != keycodes.ctrl &&
        event.keyCode != keycodes.leftWindowKey &&
        event.keyCode != keycodes.rightWindowKey &&
        !(
          (event.keyCode >= keycodes.a && event.keyCode <= keycodes.z) ||
          (event.keyCode >= keycodes.n0 && event.keyCode <= keycodes.n9) ||
          (event.keyCode >= keycodes.numPad0 && event.keyCode <= keycodes.numPad9) ||
          event.keyCode == keycodes.dash
        )
      ) {
        selectedCells.value = getCellRange(shiftDragStartCell.value || currentCell.value, currentCell.value);
      }
    }

    function handleCursorMove(event) {
      if (!isRegularCell(currentCell.value)) return null;
      let rowCount = grider.value!.rowCount;
      if (isCtrlOrCommandKey(event)) {
        switch (event.keyCode) {
          case keycodes.upArrow:
          case keycodes.pageUp:
            return moveCurrentCell(0, currentCell.value[1], event);
          case keycodes.downArrow:
          case keycodes.pageDown:
            return moveCurrentCell(rowCount - 1, currentCell.value[1], event);
          case keycodes.leftArrow:
            return moveCurrentCell(currentCell.value[0], 0, event);
          case keycodes.rightArrow:
            return moveCurrentCell(currentCell.value[0], columnSizes.value!.realCount - 1, event);
          case keycodes.home:
            return moveCurrentCell(0, 0, event);
          case keycodes.end:
            return moveCurrentCell(rowCount - 1, columnSizes.value!.realCount - 1, event);
          case keycodes.a:
            selectedCells.value = [['header', 'header']];
            event.preventDefault();
            return ['header', 'header'];
        }
      } else {
        switch (event.keyCode) {
          case keycodes.upArrow:
            if (currentCell.value[0] == 0) return focusFilterEditor(currentCell.value[1]);
            return moveCurrentCell(currentCell.value[0] - 1, currentCell.value[1], event);
          case keycodes.downArrow:
            return moveCurrentCell(currentCell.value[0] + 1, currentCell.value[1], event);
          case keycodes.enter:
            if (!grider.value?.editable) return moveCurrentCell(currentCell.value[0] + 1, currentCell.value[1], event);
            break;
          case keycodes.leftArrow:
            return moveCurrentCell(currentCell.value[0], currentCell.value[1] - 1, event);
          case keycodes.rightArrow:
            return moveCurrentCell(currentCell.value[0], currentCell.value[1] + 1, event);
          case keycodes.home:
            return moveCurrentCell(currentCell.value[0], 0, event);
          case keycodes.end:
            return moveCurrentCell(currentCell.value[0], columnSizes.value!.realCount - 1, event);
          case keycodes.pageUp:
            return moveCurrentCell(currentCell.value[0] - visibleRowCountLowerBound.value, currentCell.value[1], event);
          case keycodes.pageDown:
            return moveCurrentCell(currentCell.value[0] + visibleRowCountLowerBound.value, currentCell.value[1], event);
          case keycodes.tab: {
            return moveCurrentCellWithTabKey(event.shiftKey);
          }
        }
      }
    }

    function focusFilterEditor(columnRealIndex) {
      let modelIndex = columnSizes.value!.realToModel(columnRealIndex);
      const domFilter = domFilterControlsRef.get()[columns[modelIndex].uniqueName];
      if (domFilter) domFilter.focus()
      return ['filter', columnRealIndex]
    }

    function updateCollapsedLeftColumn() {
      void updateWidgetStyle()
      collapsedLeftColumnStore.value = !collapsedLeftColumnStore.value
    }

    function cellsToRegularCells(cells) {
      cells = flatten(
        cells.map(cell => {
          if (cell[1] == 'header') {
            return range(0, columnSizes.value!.count).map(col => [cell[0], col]);
          }
          return [cell]
        })
      );
      cells = flatten(
        cells.map(cell => {
          if (cell[0] == 'header') {
            return range(0, grider.value!.rowCount).map(row => [row, cell[1]]);
          }
          return [cell]
        })
      );
      return cells.filter(isRegularCell);
    }

    function getCellsPublished(cells) {
      const regular = cellsToRegularCells(cells);
      return regular
        .map(cell => {
          const row = cell[0];
          const rowData = grider.value!.getRowData(row);
          const column = realColumnUniqueNames[cell[1]];
          return {
            row,
            rowData,
            column,
            value: rowData && rowData[column],
            engine: display.value?.driver,
          };
        })
        .filter(x => x.column)
    }


    function showMultilineCellEditorConditional(cell) {
      if (!cell) return false
      const rowData = grider.value!.getRowData(cell[0])
      if (!rowData) return null
      const cellData = rowData[realColumnUniqueNames.value[cell[1]]]
      console.log(cellData.value, `realColumnUniqueNames`)
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
      let newFirstVisibleRowScrollIndex = firstVisibleRowScrollIndex.value
      if (deltaY > 0 && deltaX === -0) {
        newFirstVisibleRowScrollIndex += wheelRowCount.value
      } else if (deltaY < 0 && deltaX === -0) {
        newFirstVisibleRowScrollIndex -= wheelRowCount.value
      }

      let rowCount = grider.value!.rowCount
      if (newFirstVisibleRowScrollIndex + visibleRowCountLowerBound.value > rowCount) {
        newFirstVisibleRowScrollIndex = rowCount - visibleRowCountLowerBound.value + 1;
      }
      if (newFirstVisibleRowScrollIndex < 0) {
        newFirstVisibleRowScrollIndex = 0;
      }

      firstVisibleRowScrollIndex.value = newFirstVisibleRowScrollIndex;
      domVerticalScroll.value && domVerticalScroll.value.scroll(newFirstVisibleRowScrollIndex);
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

    function setCellValue(cell, value) {
      grider.value && grider.value.setCellValue(cell[0], realColumnUniqueNames.value[cell[1]], value);
    }

    function moveCurrentCell(row, col, event: Nullable<Event> = null) {
      const rowCount = grider.value!.rowCount;

      if (row < 0) row = 0;
      if (row >= rowCount) row = rowCount - 1;
      if (col < 0) col = 0;
      if (col >= columnSizes.value!.realCount) col = columnSizes.value!.realCount - 1;
      currentCell.value = [row, col];
      selectedCells.value = [[row, col]];
      scrollIntoView([row, col]);


      if (event) event.preventDefault();
      return [row, col];
    }

    function moveCurrentCellWithTabKey(isShift) {
      if (!isRegularCell(currentCell.value)) return null

      if (isShift) {
        if (currentCell[1] > 0) {
          return moveCurrentCell(currentCell[0], currentCell[1] - 1, event);
        } else {
          return moveCurrentCell(currentCell[0] - 1, columnSizes.value!.realCount - 1, event);
        }
      } else {
        if (currentCell[1] < columnSizes.value!.realCount - 1) {
          return moveCurrentCell(currentCell[0], currentCell[1] + 1, event);
        } else {
          return moveCurrentCell(currentCell[0] + 1, 0, event);
        }
      }
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

    function handleSetFormView(rowData, column) {
      if (column) {
        openReferenceForm(unref(rowData), column, conid.value, database.value);
      } else {
        openPrimaryKeyForm(unref(rowData), display.value?.baseTable, conid.value, database.value);
      }
    }

    function handleGridMouseDown(event) {
      if (event.target.closest('.buttonLike')) return;
      if (event.target.closest('.resizeHandleControl')) return;
      if (event.target.closest('.collapseButtonMarker')) return;
      if (event.target.closest('.showFormButtonMarker')) return;
      if (event.target.closest('input')) return;
      shiftDragStartCell.value = null
      // event.target.closest('table').focus();
      event.preventDefault();
      if (domFocusField.value) domFocusField.value.focus();
      const cell = cellFromEvent(event);
      if (event.button == 2) {
        if (cell && !cellIsSelected(cell[0], cell[1], selectedCells)) {
          selectedCells.value = [cell];
        }
        return;
      }
      const autofill = event.target.closest('div.autofillHandleMarker');
      if (autofill) {
        autofillDragStartCell.value = cell;
      } else {
        const oldCurrentCell = currentCell.value
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
      if (display.value && display.value.focusedColumns) display.value.focusColumns(null)
    }

    function handleGridMouseMove(event) {
      if (autofillDragStartCell.value) {
        const cell = cellFromEvent(event)
        if (isRegularCell(cell) && (cell[0] == autofillDragStartCell.value[0] || cell[1] == autofillDragStartCell.value[1])) {
          const autoFillStart = [selectedCells.value[0][0], min(selectedCells.value.map(x => x[1]))];
          // @ts-ignore
          autofillSelectedCells.value = getCellRange(autoFillStart, cell);
        }
      } else if (dragStartCell.value) {
        const cell = cellFromEvent(event);
        currentCell.value = cell;
        selectedCells.value = getCellRange(dragStartCell.value!, cell);
      }
    }

    function handleGridMouseUp(event) {
      if (dragStartCell.value) {
        const cell = cellFromEvent(event);
        currentCell.value = cell;
        selectedCells.value = getCellRange(dragStartCell.value!, cell);
        dragStartCell.value = null;
      }
      if (autofillDragStartCell.value) {
        const currentRowNumber = currentCell.value[0];
        if (isNumber(currentRowNumber)) {
          const rowIndexes = uniq((autofillSelectedCells.value || []).map(x => x[0])).filter(x => x != currentRowNumber);
          const colNames = selectedCells.value.map(cell => realColumnUniqueNames.value[cell[1]!]);
          const changeObject = pick(grider.value?.getRowData(currentRowNumber), colNames);
          grider.value?.beginUpdate();
          for (const index of rowIndexes) grider.value?.updateRow(index, changeObject);
          grider.value?.endUpdate();
        }

        autofillDragStartCell.value = null;
        autofillSelectedCells.value = [];
        selectedCells.value = autofillSelectedCells.value;
      }
    }

    return {
      container,
      domFocusField,
      domHorizontalScroll,
      domVerticalScroll,
      ...toRefs(props),
      errorMessage,
      selectedCellsInfo,
      columns,
      columnSizes,
      headerColWidth,
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
      containerWidth,
      containerHeight,
      inplaceEditorState,
      dispatchInsplaceEditor,
      firstVisibleRowScrollIndex,
      firstVisibleColumnScrollIndex,
      handleGridKeyDown,
      handleSetFormView,
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

