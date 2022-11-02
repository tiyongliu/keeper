import {Component, defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'
import LoadingDataGridCore from '/@/second/datagrid/LoadingDataGridCore'
// import eb_system_config from '/@/second/tabs/eb_system_config.json'
import credential_count from '/@/second/tabs/credential_count.json'
import {GridConfig, GridDisplay, MacroDefinition} from "/@/second/keeper-datalib";
import ChangeSetGrider from './ChangeSetGrider'
import {databaseConnectionsSqlSelectApi} from '/@/api/simpleApis'

async function loadDataPage(props, offset, limit) {
  const {display, conid, database} = props
  const select = display.getPageQuery(offset, limit)

  const response = await databaseConnectionsSqlSelectApi({
    conid: unref(conid)!,
    database: unref(database)!,
    select,
  }) as any
  return response.rows
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
    selectedCellsPublished: {
      type: Function as PropType<() => any[]>,
      default: () => []
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
  setup(props, {attrs}) {
    const grider = ref()
    const loadedRows = ref([])
    const {
      macroPreview,
      changeSetState,
      dispatchChangeSet,
      display,
      macroValues,
      selectedCellsPublished
    } = toRefs(props)

    function dataPageAvailable(props) {

    }

    async function loadRowCount() {
      return parseInt(credential_count.count)
    }

    watch(() => macroPreview.value, () => {
      if (macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!, macroPreview.value! as MacroDefinition, macroValues.value, selectedCellsPublished.value())
      }

      if (!macroPreview.value) {
        grider.value = new ChangeSetGrider(loadedRows.value, changeSetState.value, dispatchChangeSet.value, display.value!)
      }
    })

    return () => (
      <LoadingDataGridCore
        {...Object.assign({}, props, attrs)}
        loadDataPage={loadDataPage}
        loadRowCount={loadRowCount}
        display={display.value}
      />
    )
  }
})
