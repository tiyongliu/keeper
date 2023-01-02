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
  },
  emits: ['update:loadedRows', 'selectedCellsPublished'],
  setup(props, {attrs, emit}) {
    const {
      macroPreview,
      changeSetState,
      dispatchChangeSet,
      display,
      macroValues,
    } = toRefs(props)

    const grider = ref()
    const loadedRows = ref([])
    const selectedCellsPublished = ref<() => any[]>(() => [])

    function dataPageAvailable(props) {
      const {display} = props;
      const select = display.getPageQuery(0, 1);
      return !!select;
    }

    watchEffect(() => {
      if (!macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!)
      }
    })

    watchEffect(() => {
      if (macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!, macroPreview.value! as MacroDefinition, macroValues.value, selectedCellsPublished.value())
      }
    })

    function handleSelectedCellsPublished(data) {
      selectedCellsPublished.value = data
      emit('selectedCellsPublished', data)
    }

    onBeforeUnmount(() => {
      loadedRows.value = []
      grider.value = null
    })

    //onLoadedRows={rows => loadedRows.value = rows}

    //     onSelectedCellsPublished={handleSelectedCellsPublished}

    watch(() => unref(loadedRows.value), () => {
      // emit('update:loadedRows', unref(loadedRows.value))
      // console.log(loadedRows.value, `len`)
    })

    return () => (
      <LoadingDataGridCore
        {...Object.assign({}, props, attrs)}
        loadDataPage={loadDataPage}
        dataPageAvailable={dataPageAvailable}
        loadRowCount={loadRowCount}
        vModel:loadedRows={loadedRows.value}
        onSelectedCellsPublished={handleSelectedCellsPublished}
        frameSelection={!!macroPreview.value}
        grider={grider.value}
        display={display.value}
      />
    )
  }
})
