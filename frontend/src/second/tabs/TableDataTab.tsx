import {defineComponent, PropType, provide, ref, toRefs, unref} from 'vue'
import {createChangeSet, createGridCache} from '/@/second/keeper-datalib'
import ToolStripContainer from '/@/second/buttons/ToolStripContainer.vue'
import TableDataGrid from '/@/second/datagrid/TableDataGrid.vue'
import useGridConfig from '/@/second/utility/useGridConfig'
import createUndoReducer from '/@/second/utility/createUndoReducer'
import {getLocalStorage} from '/@/second/utility/storageCache'

export default defineComponent({
  name: 'TableDataTab',
  props: {
    tabid: {
      type: String as PropType<string>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    }
  },
  setup(props, {attrs}) {
    const {tabid} = toRefs(props)
    // const autoRefreshInterval = ref(getIntSettingsValue('dataGrid.defaultAutoRefreshInterval', 10, 1, 3600));
    // const autoRefreshStarted = ref(false)
    // const autoRefreshTimer = ref(null)
    // let connection = ref()
    // watch(() => conid.value, () => {
    //   useConnectionInfo({conid: unref(conid)}, connection)
    // })


    const [changeSetStore, dispatchChangeSet] = createUndoReducer(createChangeSet())

    const config = useGridConfig(tabid.value!)
    const cache = ref(createGridCache())

    function configUpdate(target) {
      config.value = target
    }

    function cacheUpdate(target) {
      cache.value = target
    }

    const collapsedLeftColumnStore = ref(getLocalStorage('dataGrid_collapsedLeftColumn', false))
    provide('collapsedLeftColumnStore', collapsedLeftColumnStore)

    return () => (
      <ToolStripContainer>
        <TableDataGrid
          {...Object.assign(props, attrs)}
          config={config.value}
          setConfig={configUpdate}
          cache={unref(cache)}
          setCache={cacheUpdate}
          focusOnVisible
          changeSetStore={changeSetStore.value}
          dispatchChangeSet={dispatchChangeSet}
        />
      </ToolStripContainer>
    )
  }
})

export const matchingProps = ['conid', 'database', 'schemaName', 'pureName'];
export const allowAddToFavorites = _ => true;
