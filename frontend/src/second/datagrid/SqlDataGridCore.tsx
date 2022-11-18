import {Component, defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'
import LoadingDataGridCore from '/@/second/datagrid/LoadingDataGridCore.vue'
// import eb_system_config from '/@/second/tabs/eb_system_config.json'
// import credential_count from '/@/second/tabs/credential_count.json'
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

  await databaseConnectionsSqlSelectApi({
    conid: unref(conid)!,
    database: unref(database)!,
    select,
  })
  return parseInt("7")
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
    macroPreview: {
      type: [String, Object] as PropType<string | Component | MacroDefinition>,
    },
    macroValues: {
      type: Object as PropType<any>
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
  },
  emits: ['selectedCellsPublished'],
  setup(props, {attrs, emit}) {
    const grider = ref()
    const loadedRows = ref([])
    const {
      macroPreview,
      changeSetState,
      dispatchChangeSet,
      display,
      macroValues,
    } = toRefs(props)

    const selectedCellsPublished = ref<() => any[]>(() => [])

    function dataPageAvailable(props) {
      const { display } = props;
      const select = display.getPageQuery(0, 1);
      return !!select;
    }

    watch(() => [macroPreview.value, ...loadedRows.value, selectedCellsPublished.value], () => {
      if (macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!, macroPreview.value! as MacroDefinition, macroValues.value, selectedCellsPublished.value())
      }

      if (!macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!)
      }
    })

    function handlerRows(rows: []) {
      loadedRows.value = rows
    }

    function handleSelectedCellsPublished(data) {
      selectedCellsPublished.value = data
      emit('selectedCellsPublished', data)
    }

    return () => (
      <LoadingDataGridCore
        {...Object.assign({}, props, attrs)}
        loadDataPage={loadDataPage}
        dataPageAvailable={dataPageAvailable}
        loadRowCount={loadRowCount}
        onLoadedRows={handlerRows}
        onSelectedCellsPublished={handleSelectedCellsPublished}
        frameSelection={!!macroPreview.value}
        grider={grider.value}
        display={display.value}
      />
    )
  }
})
