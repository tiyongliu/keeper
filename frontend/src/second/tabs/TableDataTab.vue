<template>
  <ToolStripContainer>
    <TableDataGrid
      v-bind="Object.assign({}, $props, $attrs)"
      :config="config"
      :setConfig="configUpdate"
      :cache="cache"
      :setCache="cacheUpdate"
      focusOnVisible
      :changeSetState="changeSetStore"
      :dispatchChangeSet="dispatchChangeSet"
    />
    <template #toolstrip>
      <ToolStripCommandSplitButton
        :buttonLabel="autoRefreshStarted ? `Refresh (every ${autoRefreshInterval}s)` : null"
        :commands="['dataGrid.refresh', ...createAutoRefreshMenu()]"
        hideDisabled
      />

      <!-- <ToolStripCommandButton command="dataGrid.refresh" hideDisabled />
      <ToolStripCommandButton command="dataForm.refresh" hideDisabled /> -->

      <ToolStripCommandButton command="dataForm.goToFirst" hideDisabled />
      <ToolStripCommandButton command="dataForm.goToPrevious" hideDisabled />
      <ToolStripCommandButton command="dataForm.goToNext" hideDisabled />
      <ToolStripCommandButton command="dataForm.goToLast" hideDisabled />

      <ToolStripCommandButton command="tableData.save" />
      <ToolStripCommandButton command="dataGrid.insertNewRow" hideDisabled />
      <ToolStripCommandButton command="dataGrid.deleteSelectedRows" hideDisabled />
      <ToolStripCommandButton command="dataGrid.switchToForm" hideDisabled />
      <ToolStripCommandButton command="dataGrid.switchToTable" hideDisabled />
    </template>
  </ToolStripContainer>
</template>

<script lang="ts">
import {defineComponent, onBeforeUnmount, PropType, provide, ref, toRefs, unref, watch} from 'vue'
import {getLocalStorage, setLocalStorage} from '/@/second/utility/storageCache'
import {createChangeSet, createGridCache} from '/@/second/keeper-datalib'
import ToolStripContainer from '/@/second/buttons/ToolStripContainer.vue'
import TableDataGrid from '/@/second/datagrid/TableDataGrid.vue'
import ToolStripCommandSplitButton from '/@/second/buttons/ToolStripCommandSplitButton.vue'
import ToolStripCommandButton from '/@/second/buttons/ToolStripCommandButton.vue'
import createUndoReducer from '/@/second/utility/createUndoReducer'
import useGridConfig from '/@/second/utility/useGridConfig'

export const matchingProps = ['conid', 'database', 'schemaName', 'pureName']
export const allowAddToFavorites = _ => true
const INTERVALS = [5, 10, 15, 13, 60]
export default defineComponent({
  name: 'TableDataTab',
  components: {
    ToolStripContainer,
    TableDataGrid,
    ToolStripCommandButton,
    ToolStripCommandSplitButton
  },
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
  setup(props) {
    const {tabid} = toRefs(props)

    let autoRefreshTimer: number | null = null


    const autoRefreshInterval = ref(0)
    const autoRefreshStarted = ref(false)


    // const autoRefreshInterval = ref(getIntSettingsValue('dataGrid.defaultAutoRefreshInterval', 10, 1, 3600));
    // const autoRefreshTimer = ref(null)
    // let connection = ref()
    // watch(() => conid.value, () => {
    //   useConnectionInfo({conid: unref(conid)}, connection)
    // })

    const [changeSetStore, dispatchChangeSet] = createUndoReducer(createChangeSet())

    const config = useGridConfig(tabid.value!)
    const cache = ref(createGridCache())

    function configUpdate(updater) {
      if (updater) config.value = updater(config.value)
    }

    function cacheUpdate(updater) {
      if (updater) cache.value = updater(cache.value)
    }

    const collapsedLeftColumnStore = ref(getLocalStorage('dataGrid_collapsedLeftColumn', false))
    provide('collapsedLeftColumnStore', collapsedLeftColumnStore)

    watch(() => collapsedLeftColumnStore.value, (_, newValue) => {
      setLocalStorage('dataGrid_collapsedLeftColumn', unref(newValue))
    })

    function closeRefreshTimer() {
      if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer)
        autoRefreshTimer = null
      }
    }

    onBeforeUnmount(() => closeRefreshTimer())

    function createAutoRefreshMenu() {
      return [
        { divider: true },
        { command: 'tableData.stopAutoRefresh', hideDisabled: true },
        { command: 'tableData.startAutoRefresh', hideDisabled: true },
        'tableData.setAutoRefresh.1',
        ...INTERVALS.map(seconds => ({ command: `tableData.setAutoRefresh.${seconds}`, text: `...${seconds} seconds` })),
      ]
    }

    return {
      config,
      configUpdate,
      cache,
      cacheUpdate,
      changeSetStore,
      dispatchChangeSet,
      createAutoRefreshMenu,
      autoRefreshStarted,
      autoRefreshInterval
    }
  }
})
</script>
