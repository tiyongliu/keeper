import {
  computed,
  defineComponent,
  onMounted,
  PropType,
  provide,
  ref,
  toRefs,
  unref,
  watch
} from 'vue'
import {storeToRefs} from 'pinia'
import ToolStripContainer from '/@/second/buttons/ToolStripContainer.vue'
import DataGrid from '/@/second/datagrid/DataGrid.vue'
import CollectionDataGridCore from '/@/second/datagrid/CollectionDataGridCore'
import CollectionJsonView from '/@/second/jsonview/CollectionJsonView.vue'
import {
  CollectionGridDisplay,
  createChangeSet,
  createGridCache,
  runMacroOnChangeSet
} from '/@/second/keeper-datalib'
import useGridConfig from '/@/second/utility/useGridConfig'
import {useCollectionInfo, useConnectionInfo} from '/@/api/bridge'
import {findEngineDriver} from '/@/second/keeper-tools'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import createUndoReducer from '/@/second/utility/createUndoReducer'
import {getLocalStorage, setLocalStorage} from '/@/second/utility/storageCache'

export default defineComponent({
  name: 'CollectionDataTab',
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
    const {tabid, conid, database, schemaName, pureName} = toRefs(props)
    const connection = ref()
    const collectionInfo = ref()
    const loadedRows = ref([])

    const bootstrap = useBootstrapStore()
    const {extensions} = storeToRefs(bootstrap)

    const config = useGridConfig(tabid.value!)
    const cache = ref(createGridCache())

    function configUpdate(updater) {
      if (updater) config.value = updater(config.value)
    }

    function cacheUpdate(updater) {
      if (updater) cache.value = updater(cache.value)
    }

    const [changeSetStore, dispatchChangeSet] = createUndoReducer(createChangeSet())

    watch(() => conid.value, () => {
      useConnectionInfo({conid: unref(conid)}, connection)
    }, {immediate: true})

    watch(() => [conid.value, database.value, schemaName.value, pureName.value], () => {
      useCollectionInfo({
        conid: conid.value,
        database: database.value,
        schemaName: schemaName.value,
        pureName: pureName.value
      }, collectionInfo)
    }, {immediate: true})

    const display = computed(() => collectionInfo.value && connection.value
      ? new CollectionGridDisplay(
        collectionInfo.value,
        findEngineDriver(connection.value, extensions.value!),
        //@ts-ignore
        config.value,
        configUpdate,
        cache.value,
        cacheUpdate,
        loadedRows.value,
        unref(changeSetStore)?.value,
        connection.value?.isReadOnly,
      )
      : null
    )

    onMounted(() => {

    })

    function handleRunMacro(macro, params, cells) {
      const newChangeSet = runMacroOnChangeSet(macro, params, cells, unref(changeSetStore)?.value, display.value!)
      if (newChangeSet) {
        dispatchChangeSet({type: 'set', value: newChangeSet});
      }
    }

    const collapsedLeftColumnStore = ref(getLocalStorage('collection_collapsedLeftColumn', false))
    provide('collapsedLeftColumnStore', collapsedLeftColumnStore)

    watch(() => collapsedLeftColumnStore.value, (_, newValue) => {
      setLocalStorage('collection_collapsedLeftColumn', unref(newValue))
    })

    return () => (
      <ToolStripContainer>
        <DataGrid
          vModel:loadedRows={loadedRows.value}
          {...Object.assign({}, props, attrs)}
          config={config.value}
          setConfig={configUpdate}
          cache={cache.value}
          setCache={cacheUpdate}
          focusOnVisible
          display={display.value!}
          changeSetState={changeSetStore.value}
          dispatchChangeSet={dispatchChangeSet}
          gridCoreComponent={CollectionDataGridCore}
          jsonViewComponent={CollectionJsonView}
          isDynamicStructure
          showMacros
          macroCondition={macro => macro.type == 'transformValue'}
          runMacro={handleRunMacro}
        />
      </ToolStripContainer>
    )
  },
})

export const matchingProps = ['conid', 'database', 'schemaName', 'pureName'];
export const allowAddToFavorites = _ => true
