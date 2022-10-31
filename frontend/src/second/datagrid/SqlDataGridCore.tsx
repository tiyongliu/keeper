import {omit} from 'lodash-es'
import {Component, defineComponent, onMounted, PropType, ref, toRefs, unref, watch} from 'vue'
import LoadingDataGridCore from '/@/second/datagrid/LoadingDataGridCore'
import eb_system_config from '/@/second/tabs/eb_system_config.json'
import credential_count from '/@/second/tabs/credential_count.json'
import {GridConfig, GridDisplay, MacroDefinition} from "/@/second/keeper-datalib";
import ChangeSetGrider from './ChangeSetGrider'
import {databaseConnectionsSqlSelectApi} from '/@/api/simpleApis'

async function loadDataPage(props, offset, limit) {
  const { display } = props


  console.log(`-----------------------props`, )
  console.log(`-----------------------offset`, offset)
  console.log(`-----------------------limit`, limit)
  console.log(`-----------------------display`, display)
  // const response = await databaseConnectionsSqlSelectApi({
  //   conid: unref(conid)!,
  //   database: unref(database)!,
  //   select,
  // });


  return eb_system_config.rows
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
      conid,
      database,
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

    const fullProps = Object.assign({}, props, attrs)

    onMounted(() => {
      setTimeout(() => {
        console.log(`-----------??????????fullProps`, fullProps)
      }, 5000)
    })

    return () => (
      <LoadingDataGridCore
        {...fullProps}
        loadDataPage={loadDataPage}
        loadRowCount={loadRowCount}
        display={display.value}
      />
    )
  }
})
