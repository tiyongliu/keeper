import {Component, defineComponent, PropType, ref, toRefs, watch} from 'vue'
import LoadingDataGridCore from '/@/second/datagrid/LoadingDataGridCore'
import eb_system_config from '/@/second/tabs/eb_system_config.json'
import credential_count from '/@/second/tabs/credential_count.json'
import {GridConfig, TableGridDisplay} from "/@/second/keeper-datalib";
import ChangeSetGrider from './ChangeSetGrider'

export default defineComponent({
  name: 'SqlDataGridCore',
  props: {
    conid: {
      type: String as PropType<string>
    },
    display: {
      type: Object as PropType<TableGridDisplay>
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
      type: [String, Object] as PropType<string | Component>,
    },
    selectedCellsPublished: {
      type: Function as PropType<() => any[]>,
      default: () => []
    }
  },
  setup(props, {attrs}) {
    const grider = ref()
    const {macroPreview} = toRefs(props)

    async function loadDataPage() {
      return eb_system_config.rows
    }

    function dataPageAvailable(props) {

    }
    
    async function loadRowCount() {
      return parseInt(credential_count.count)
    }


    watch(() => macroPreview.value, () => {
      if (macroPreview.value) {

      }
    })

    const fullProps = Object.assign(props, attrs)
    return () => (
      <LoadingDataGridCore
        {...fullProps}
        loadDataPage={loadDataPage}
        loadRowCount={loadRowCount}
      />
    )
  }
})
