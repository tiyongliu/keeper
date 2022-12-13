<template>
  <div class="outer">
    <div class="wrapper" ref="container">
      <table v-for="(chunk, chunkIndex) in columnChunks" :key="chunkIndex">
        <tr v-for="(col, rowIndex) in chunk" :key="rowIndex">
          <td
            class="header-cell"
            :data-row="`${rowIndex}`"
            :data-col="`${chunkIndex * 2}`"
            :style="rowHeight > 1 ? `height: ${rowHeight}px` : undefined"
            :class="[formDisplay.config.formColumnFilterText &&
                filterName(formDisplay.config.formColumnFilterText, col.columnName) && 'columnFiltered',
                (currentCell[0] == rowIndex && currentCell[1] == chunkIndex * 2) && 'isSelected'
            ]">
            <div class="header-cell-inner">
              <FontIcon
                v-if="col.foreignKey"
                :icon="plusExpandIcon(formDisplay.isExpandedColumn(col.uniqueName))"
                @click.stop="formDisplay.toggleExpandedColumn(col.uniqueName)"/>
              <FontIcon v-else icon="icon invisible-box"/>
              <span :style="`margin-left: ${(col.uniquePath.length - 1) * 20}px`"></span>
              <ColumnLabel
                v-bind="{...col}"
                :headerText="col.columnName"
                showDataType
                :conid="conid"
                :database="database"/>
            </div>
          </td>
          <DataGridCell
            :maxWidth="(wrapperWidth * 2) / 3"
            :minWidth="200"
            :rowIndex="rowIndex"
            :col="col"
            :rowData="rowData"
            :colIndex="chunkIndex * 2 + 1"
            :isSelected="currentCell[0] == rowIndex && currentCell[1] == chunkIndex * 2 + 1"
            :isModifiedCell="(rowStatus && rowStatus.modifiedFields) && rowStatus.modifiedFields.has(col.uniqueName)"
            :allowHintField="!((rowStatus && rowStatus.modifiedFields) && rowStatus.modifiedFields.has(col.uniqueName))"
            v-model:domCell="domCells[`${rowIndex},${chunkIndex * 2 + 1}`]"
            :setFormView="handleSetFormView"
            :showSlot="!rowData ||
                  (inplaceEditorState.cell &&
                  rowIndex == inplaceEditorState.cell[0] &&
                  chunkIndex * 2 + 1 == inplaceEditorState.cell[1])"
            :isCurrentCell="currentCell[0] == rowIndex && currentCell[1] == chunkIndex * 2 + 1"
            :dictionaryLookup="() => handleLookup(col)"
          >
            <InplaceEditor
              v-if="rowData &&
              inplaceEditorState.cell &&
              rowIndex == inplaceEditorState.cell[0] &&
              chunkIndex * 2 + 1 == inplaceEditorState.cell[1]"
              fillParent
              :width="getCellWidth(rowIndex, chunkIndex * 2 + 1)"
              :inplaceEditorState="inplaceEditorState"
              :dispatchInsplaceEditor="dispatchInsplaceEditor"
              :cellValue="rowData[col.uniqueName]"
              @setValue="value => former.setCellValue(col.uniqueName, value)"/>
          </DataGridCell>
        </tr>
      </table>
      <input
        type="text"
        class="focus-field"
        @keydown="handleKeyDown"
        @copy="copyToClipboard"
      />
    </div>
    <div v-if="rowCountInfo" class="row-count-label">{{ rowCountInfo }}</div>
  </div>

  <LoadingInfo v-if="isLoading" wrapper message="Loading data"/>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  nextTick,
  onMounted,
  PropType,
  reactive,
  ref,
  toRefs,
  unref,
  watch
} from 'vue'
import {chunk} from 'lodash-es'
import {Input} from 'ant-design-vue'
import {filterName} from '/@/second/keeper-tools'
import {GridConfig, TableFormViewDisplay} from '/@/second/keeper-datalib'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import DataGridCell from '/@/second/datagrid/DataGridCell.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import InplaceEditor from '/@/second/datagrid/InplaceEditor.vue'
import {plusExpandIcon} from '/@/second/icons/expandIcons'
import createReducer from '/@/second/utility/createReducer'
import ChangeSetFormer from './ChangeSetFormer'
import {extractRowCopiedValue, copyTextToClipboard} from '/@/second/utility/clipboard'
export default defineComponent({
  name: 'FormView',
  components: {
    FontIcon,
    ColumnLabel,
    DataGridCell,
    LoadingInfo,
    InplaceEditor,
    [Input.name]: Input
  },
  props: {
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    config: {
      type: Object as PropType<GridConfig>,
    },
    focusOnVisible: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    setConfig: {
      type: Function as PropType<(target: any) => void>
    },
    allRowCount: {
      type: Number as PropType<Nullable<number>>
    },
    rowCountBefore: {
      type: Number as PropType<Nullable<number>>
    },
    isLoading: {
      type: Boolean as PropType<boolean>
    },
    former: {
      type: Object as PropType<Nullable<ChangeSetFormer>>
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    navigate: {
      type: Function as PropType<(command: Promise<'begin' | 'previous' | 'next' | 'end'>) => void>
    }
  },
  setup(props) {
    const {former, rowCountBefore, allRowCount, focusOnVisible, formDisplay} = toRefs(props)
    const container = ref<Nullable<HTMLElement>>(null)
    const domCells = ref<{ [key in string]: HTMLElement }>({})

    const wrapperWidth = ref(container.value ? container.value.clientWidth : 0)
    const wrapperHeight = ref(container.value ? container.value.clientWidth : 0)
    const rowHeight = computed(() => 25) //todo  $: rowHeight = $dataGridRowHeight;

    let currentCell = reactive([0, 0])

    const tabVisible = inject('tabVisible')

    watch(() => [unref(tabVisible), focusOnVisible.value], () => {
      if (unref(tabVisible) && focusOnVisible.value) {
        updateWidgetStyle()
      }
    })

    const rowData = computed(() => former.value ? former.value?.rowData : null)
    const rowStatus = computed(() => former.value ? former.value?.rowStatus : null)

    const rowCount = computed(() => Math.floor((wrapperHeight.value - 22) / (rowHeight.value + 2)))

    const columnChunks = computed(() => formDisplay.value
      ? chunk(formDisplay.value.columns, rowCount.value) as any[][]
      : [])

    const rowCountInfo = computed(() => getRowCountInfo(rowCountBefore.value, allRowCount.value))

    function getRowCountInfo(rowCountBefore, allRowCount) {
      if (rowData.value == null) return 'No data'
      if (allRowCount == null || rowCountBefore == null) return 'Loading row count...'
      return `Row: ${(rowCountBefore + 1).toLocaleString()} / ${allRowCount.toLocaleString()}`
    }

    function updateWidgetStyle() {
      nextTick(() => {
        if (container.value && container.value!.clientWidth) wrapperWidth.value = container.value!.clientWidth
        if (container.value && container.value!.clientHeight) wrapperHeight.value = container.value!.clientHeight
      })
    }

    onMounted(() => {
      updateWidgetStyle()
    })

    function copyToClipboard() {
      const column = getCellColumn(currentCell);
      if (!column) return;
      const text = currentCell[1] % 2 == 1 ? extractRowCopiedValue(rowData, column.uniqueName) : column.columnName;
      copyTextToClipboard(text);
    }

    const scrollIntoView = cell => {
      const element = domCells[`${cell[0]},${cell[1]}`];
      if (element) element.scrollIntoView();
    };

    const moveCurrentCell = (row, col) => {
      if (row < 0) row = 0;
      if (col < 0) col = 0;
      if (col >= columnChunks.value.length * 2) col = columnChunks.value.length * 2 - 1;
      const chunk = columnChunks.value[Math.floor(col / 2)];
      if (chunk && row >= chunk.length) row = chunk.length - 1;
      currentCell = [row, col];
      scrollIntoView(currentCell);
    }

    function getCellColumn(cell) {
      const chunk = columnChunks.value[Math.floor(cell[1] / 2)]
      if (!chunk) return;
      return chunk[cell[0]]
    }

    const getCellWidth = (row, col) => {
      const element = domCells.value[`${row},${col}`];
      if (element) return element.getBoundingClientRect().width;
      return 100;
    }

    const [inplaceEditorState, dispatchInsplaceEditor] = createReducer((state, action) => {
      switch (action.type) {
        case 'show': {
          if (former.value && !former.value.editable) return {}
          const column = getCellColumn(unref(action).cell)
          if (!unref(column)) return state
          if (unref(column).uniquePath.length > 1) return state

          // if (!grider.editable) return {};
          return {
            cell: unref(action).cell,
            text: unref(action).text,
            selectAll: unref(action).selectAll,
          };
        }
        case 'close': {
          const [row, col] = currentCell || []
          if (action.mode == 'enter' && row) setTimeout(() => moveCurrentCell(row + 1, col), 0)
          return {}
        }
      }
      return {}
    }, {})

    function handleSetFormView(rowData, column) {

    }

    function handleLookup(col) {
      console.log(col, `col?col`)
    }

    function handleKeyDown(event) {

    }
    return {
      container,
      domCells,
      ...toRefs(props),
      wrapperWidth,
      currentCell,
      rowData,
      rowStatus,
      rowHeight,
      columnChunks,
      rowCountInfo,
      inplaceEditorState,
      dispatchInsplaceEditor,
      getCellWidth,
      filterName,
      plusExpandIcon,
      handleLookup,
      handleSetFormView,
      handleKeyDown,
      copyToClipboard
    }
  }
})
</script>

<style scoped>
table {
  border-collapse: collapse;
  outline: none;
}

.outer {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
}

.wrapper {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
  display: flex;
  overflow-x: scroll;
  align-items: flex-start;
}

tr {
  background-color: var(--theme-bg-0);
}

tr:nth-child(6n + 3) {
  background-color: var(--theme-bg-1);
}

tr:nth-child(6n + 6) {
  background-color: var(--theme-bg-alt);
}

.header-cell {
  border: 1px solid var(--theme-border);
  text-align: left;
  padding: 0;
  margin: 0;
  background-color: var(--theme-bg-1);
  overflow: hidden;
}

.header-cell.isSelected {
  background: var(--theme-bg-selected);
}

.header-cell-inner {
  display: flex;
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

.columnFiltered {
  background: var(--theme-bg-green);
}
</style>
