<template>
  <VerticalSplitter :isSplitter="!!reference">
    <template #1>
      <DataGrid
        v-bind="Object.assign({}, $props, $attrs)"
        :gridCoreComponent="SqlDataGridCore"
        :formViewComponent="SqlFormView"
        :display="display"
        :formDisplay="formDisplay"
        showReferences
        showMacros
        :runMacro="handleRunMacro"
        :macroCondition="macro => macro.type == 'transformValue'"
        :multipleGridsOnTab="multipleGridsOnTab || !!reference"
        allowDefineVirtualReferences
      />
    </template>
    <template #2>
      <div class="reference-container">
        <ReferenceHeader :reference="reference" @close="handleCloseReference"/>
        <div class="detail">
          <TableDataGrid
            v-bind="Object.assign({}, $props, $attrs)"
            :pureName="reference.pureName"
            :schemaName="reference.schemaName"
            :config="childConfig"
            :setConfig="setChildConfig"
            :cache="childCache"
            :setCache="childCacheUpdate"
            :masterLoadedTime="myLoadedTime"
            isDetailView
            multipleGridsOnTab
          />
        </div>
      </div>
    </template>
  </VerticalSplitter>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  toRefs,
  unref,
  watch,
} from 'vue'
import {storeToRefs} from 'pinia'
import {fromPairs, isFunction} from 'lodash-es'
import stableStringify from 'json-stable-stringify'
import VerticalSplitter from '/@/second/elements/VerticalSplitter.vue'
import DataGrid from '/@/second/datagrid/DataGrid.vue'
import ReferenceHeader from '/@/second/datagrid/ReferenceHeader.vue'
import SqlDataGridCore from '/@/second/datagrid/SqlDataGridCore'
import SqlFormView from '/@/second/formview/SqlFormView'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {
  useConnectionInfo,
  useConnectionList,
  useDatabaseInfo,
  useDatabaseServerVersion
} from '/@/api/bridge'
import {getBoolSettingsValue} from '/@/second/settings/settingsTools'
import {getDictionaryDescription} from '/@/second/utility/dictionaryDescriptionTools'
import {
  createGridCache,
  FormViewDisplay,
  GridCache,
  GridConfig,
  GridDisplay,
  runMacroOnChangeSet,
  TableFormViewDisplay,
  TableGridDisplay
} from '/@/second/keeper-datalib'
import {extendDatabaseInfoFromApps, findEngineDriver} from '/@/second/keeper-tools'
import {getFilterValueExpression} from '/@/second/keeper-filterparser'
import {DatabaseInfo, ExtensionsDirectory} from '/@/second/keeper-types'

export default defineComponent({
  name: 'TableDataGrid',
  props: {
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
    },
    setConfig: {
      type: Function as PropType<(target: any) => void>
    },
    setCache: {
      type: Function as PropType<(target: any) => void>
    },
    config: {
      type: Object as PropType<GridConfig>,
    },
    cache: {
      type: Object as PropType<{ refreshTime: number }>
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    changeSetStore: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
    multipleGridsOnTab: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  components: {
    VerticalSplitter,
    DataGrid,
    ReferenceHeader,
    SqlDataGridCore,
    SqlFormView
  },
  setup(props) {
    const bootstrap = useBootstrapStore()
    const {extensions} = storeToRefs(bootstrap)
    const {
      config,
      setConfig,
      cache,
      setCache,
      conid,
      database,
      schemaName,
      pureName,
      changeSetState,
      dispatchChangeSet,
      multipleGridsOnTab
    } = toRefs(props)

    const reference = computed(() => config.value!.reference)
    const childConfig = computed(() => config.value!.childConfig)

    const myLoadedTime = ref(0)
    let connection = ref()
    let dbinfo = ref()
    let serverVersion = ref()
    let apps = ref([])
    let extendedDbInfo = ref()
    let connections = ref()

    watch(() => [conid.value, database.value], () => {
      useConnectionInfo({conid: unref(conid)}, connection)
      useDatabaseInfo<DatabaseInfo>({conid: unref(conid), database: unref(database)}, dbinfo)
      useDatabaseServerVersion({conid: unref(conid), database: unref(database)}, serverVersion)
    }, {immediate: true})

    watch(() => dbinfo.value, () => {
      extendedDbInfo.value = extendDatabaseInfoFromApps(dbinfo.value, apps.value)
    })

    const display = computed(() => connection.value && serverVersion.value ? new TableGridDisplay(
      {schemaName: schemaName.value, pureName: pureName.value!},
      findEngineDriver(connection.value, <ExtensionsDirectory>extensions.value!),
      config.value!,
      setConfig.value as (changeFunc: (config: GridConfig) => GridConfig) => void,
      cache.value!,
      setCache.value as (changeFunc: (cache: GridCache) => GridCache) => void,
      extendedDbInfo.value,
      {showHintColumns: getBoolSettingsValue('dataGrid.showHintColumns', true)},
      serverVersion.value,
      table => getDictionaryDescription(table, conid.value!, database.value!, apps.value, connections.value) as any
    ) as GridDisplay : null)

    const formDisplay = computed(() => connection.value && serverVersion.value ? new TableFormViewDisplay(
      {schemaName: schemaName.value, pureName: pureName.value!},
      findEngineDriver(connection.value, <ExtensionsDirectory>extensions.value!),
      config.value!,
      setConfig.value as (changeFunc: (config: GridConfig) => GridConfig) => void,
      cache.value!,
      setCache.value as (changeFunc: (cache: GridCache) => GridCache) => void,
      extendedDbInfo.value,
      {showHintColumns: getBoolSettingsValue('dataGrid.showHintColumns', true)},
      serverVersion.value,
      table => getDictionaryDescription(table, conid.value!, database.value!, apps.value, connections.value) as any
    ) as FormViewDisplay : null)

    onMounted(() => useConnectionList(connections))
    onBeforeUnmount(() => connections.value = null)

    const childCache = ref(createGridCache())
    const childCacheUpdate = target => childCache.value = target

    const setChildConfig = (value, reference: undefined | null = undefined) => {
      if (isFunction(value)) {
        setConfig.value!(x => ({
          ...x,
          childConfig: value(x.childConfig),
        }));
      } else {
        setConfig.value!(x => ({
          ...x,
          childConfig: value,
          reference: reference === undefined ? x.reference : reference,
        }));
      }
    }

    const handleReferenceSourceChanged = (selectedRows, loadedTime) => {
      myLoadedTime.value = loadedTime
      if (!reference.value) return

      const filtersBase = display.value && display.value.isGrouped ? config.value!.filters : childConfig.value!.filters

      const filters = {
        ...filtersBase,
        ...fromPairs(
          reference.value!.columns.map(col => [
            col.refName,
            selectedRows.map(x => getFilterValueExpression(x[col.baseName], (col as any).dataType)).join(', '),
          ])
        ),
      }

      if (stableStringify(filters) != stableStringify(childConfig.value!.filters)) {
        setChildConfig(cfg => ({
          ...cfg,
          filters,
        }))
        childCache.value = {
          ...childCache.value,
          refreshTime: new Date().getTime(),
        }
      }
    }

    const handleCloseReference = () => {
      display.value && display.value.clearGrouping()
      setChildConfig(null, null)
    }

    function handleRunMacro(macro, params, cells) {
      const newChangeSet = runMacroOnChangeSet(macro, params, cells, unref(changeSetState)?.value, unref(display)!)
      if (newChangeSet) {
        dispatchChangeSet.value!({type: 'set', value: newChangeSet});
      }
    }

    return {
      SqlDataGridCore,
      SqlFormView,
      handleReferenceSourceChanged,
      handleCloseReference,
      display,
      formDisplay,
      handleRunMacro,
      reference,
      multipleGridsOnTab,
      setChildConfig,
      childCacheUpdate,
      childCache: childCache.value,
      myLoadedTime: myLoadedTime.value,
      childConfig: childConfig.value,
    }
  }
})
</script>

<style scoped>
.reference-container {
  position: absolute;
  display: flex;
  flex-direction: column;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.detail {
  position: relative;
  flex: 1;
}
</style>
