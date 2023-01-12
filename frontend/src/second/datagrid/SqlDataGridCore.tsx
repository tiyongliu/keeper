import {
  Component,
  defineComponent,
  onBeforeUnmount,
  PropType,
  ref,
  toRefs,
  unref,
  watch,
  watchEffect
} from 'vue'
import LoadingDataGridCore from '/@/second/datagrid/LoadingDataGridCore'
import {GridConfig, GridDisplay, MacroDefinition} from "/@/second/keeper-datalib";
import ChangeSetGrider from './ChangeSetGrider'
import {databaseConnectionsSqlSelectApi} from '/@/api/simpleApis'

//这个要写活，查看node源码是怎么写的。
async function loadDataPage(props, offset, limit) {
  const {display, conid, database} = props
  const select = display.getPageQuery(offset, limit)
  const response = await databaseConnectionsSqlSelectApi({
    conid: unref(conid)!,
    database: unref(database)!,
    select,
  }) as any
  if (response.errorMessage) return response;
  return response.rows
}

async function loadRowCount(props) {
  const {display, conid, database} = props

  const select = display.getCountQuery()

  const response = await databaseConnectionsSqlSelectApi<{ msgtype: string; rows: { count: number }[] }>({
    conid: unref(conid)!,
    database: unref(database)!,
    select,
  })
  return response.rows[0].count
}

export default defineComponent({
  name: 'SqlDataGridCore',
  props: {
    conid: {
      type: String as PropType<string>
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    },
    config: {
      type: Object as PropType<GridConfig>,
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
    macroPreview: {
      type: [String, Object] as PropType<string | Component | MacroDefinition>,
    },
    macroValues: {
      type: Object as PropType<any>
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: []
    },
    selectedCellsPublished: {
      type: Function as PropType<() => []>,
      default: () => []
    }
  },
  emits: ['update:loadedRows', 'update:selectedCellsPublished'],
  setup(props, {attrs, emit}) {
    const {
      macroPreview,
      changeSetState,
      dispatchChangeSet,
      display,
      macroValues,
      loadedRows,
      selectedCellsPublished,
    } = toRefs(props)

    const grider = ref()
    const loadedRowsRw = ref(loadedRows.value)
    const selectedCellsPublishedRw = ref(selectedCellsPublished.value)

    function dataPageAvailable(props) {
      const {display} = props;
      const select = display.getPageQuery(0, 1);
      return !!select;
    }

    watchEffect(() => {
      if (!macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRowsRw.value, changeSetState.value, dispatchChangeSet.value, display.value!)
      }
    })

    watchEffect(() => {
      if (macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRowsRw.value, changeSetState.value, dispatchChangeSet.value, display.value!, macroPreview.value! as MacroDefinition, macroValues.value, selectedCellsPublished.value())
      }
    })

    onBeforeUnmount(() => {
      loadedRowsRw.value = []
      grider.value = null
    })

    watch(() => [...loadedRowsRw.value], () => {
      emit('update:loadedRows', unref(loadedRowsRw.value))
    })

    watch(() => selectedCellsPublishedRw.value, () => {
      emit('update:selectedCellsPublished', selectedCellsPublishedRw.value)
    })

    return () => (
      <LoadingDataGridCore
        {...Object.assign({}, props, attrs)}
        loadDataPage={loadDataPage}
        dataPageAvailable={dataPageAvailable}
        loadRowCount={loadRowCount}
        vModel:loadedRows={loadedRowsRw.value}
        vModel:selectedCellsPublished={selectedCellsPublishedRw.value}
        frameSelection={!!macroPreview.value}
        grider={grider.value}
        display={display.value}
      />
    )
  }
})
