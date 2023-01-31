import {defineComponent, PropType, ref, toRefs, watchEffect, nextTick, watch} from 'vue'
import {getIntSettingsValue} from '/@/second/settings/settingsTools'
import createRef from '/@/second/utility/createRef'
import DataGridCore from './DataGridCore.vue'
import {GridDisplay} from '/@/second/keeper-datalib'
import {isFunction,} from '/@/utils/is'
import Grider from '/@/second/datagrid/Grider'

export default defineComponent({
  name: 'LoadingDataGridCore',
  props: {
    loadDataPage: {
      type: Function as PropType<(props: any, offset: any, limit: any) => Promise<any>>,
    },
    dataPageAvailable: {
      type: Function as PropType<(props: any) => boolean>,
    },
    loadRowCount: {
      type: Function as PropType<(props: any) => Promise<number>>,
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: () => []
    },
    grider: {
      type: Object as PropType<Grider>
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    masterLoadedTime: {
      type: Number as PropType<number>,
    },
    rowCountLoaded: {
      type: Number as PropType<number>,
    },
    selectedCellsPublished: {
      type: Function as PropType<() => []>,
      default: () => []
    }
  },
  emits: ['update:selectedCellsPublished', 'update:loadedRows'],
  setup(props, {attrs, emit}) {
    const {
      loadDataPage,
      loadedRows,
      selectedCellsPublished,
      loadRowCount,
      display,
      dataPageAvailable,
      grider,
      masterLoadedTime,
      rowCountLoaded
    } = toRefs(props)

    const selectedCellsPublishedRW = ref(selectedCellsPublished.value)
    const isLoadedAll = ref<boolean>(false)
    const allRowCount = ref<Nullable<number>>(null)
    const errorMessage = ref<Nullable<string>>(null)
    const loadedTime = ref(new Date().getTime())
    const domGrid = ref<Nullable<HTMLElement>>(null)
    const loadNextDataRef = createRef<boolean>(false)
    const loadedTimeRef = createRef<number | boolean | null>(null)
    const isLoading = ref(false)

    const handleLoadRowCount = async () => {
      allRowCount.value = await loadRowCount.value!(Object.assign({}, props, attrs))
    }

    async function loadNextData() {
      if (isLoading.value) return
      loadedTimeRef.set(false)
      isLoading.value = true

      const loadStart = new Date().getTime()
      await nextTick()

      loadedTimeRef.set(loadStart)

      const nextRows = await loadDataPage.value!(
        Object.assign({}, props, attrs),
        loadedRows.value.length,
        getIntSettingsValue('dataGrid.pageSize', 100, 5, 1000)
      )

      if (loadedTimeRef.get() !== loadStart) {
        return
      }

      isLoading.value = false

      if (nextRows.errorMessage) {
        errorMessage.value = nextRows.errorMessage
      } else {
        if (allRowCount.value == null) await handleLoadRowCount()
        const loadedRowsRW = [...loadedRows.value, ...nextRows]
        emit('update:loadedRows', loadedRowsRW)
        isLoadedAll.value = nextRows.length === 0
      }

      if (loadNextDataRef.get()) {
        loadNextData()
      }
    }

    function handleLoadNextData() {
      if (!isLoadedAll.value && !errorMessage.value && (grider.value && !grider.value.disableLoadNextPage)) {
        if (isFunction(dataPageAvailable.value) && dataPageAvailable.value!(Object.assign({}, props, attrs))) {
          void loadNextData()
        }
      }
    }

    function reload() {
      allRowCount.value = null
      isLoading.value = false
      emit('update:loadedRows', [])
      isLoadedAll.value = false
      loadedTime.value = new Date().getTime()
      errorMessage.value = null
      loadNextDataRef.set(false)
    }

    watchEffect(() => {
      if ((display.value! && display.value?.cache?.refreshTime) > loadedTime.value) {
        reload()
      }
    })

    watchEffect(() => {
      if (masterLoadedTime.value && masterLoadedTime.value > loadedTime.value && display.value) {
        display.value.reload()
      }
    })

    watch(() => selectedCellsPublishedRW.value, () => {
      emit('update:selectedCellsPublished', selectedCellsPublishedRW.value)
    })

    return () => <DataGridCore
      {...Object.assign({}, props, attrs)}
      ref={domGrid}
      vModel:selectedCellsPublished={selectedCellsPublishedRW.value}
      onLoadNextData={handleLoadNextData}
      errorMessage={errorMessage.value}
      isLoading={isLoading.value}
      allRowCount={(rowCountLoaded.value || allRowCount.value)!}
      isLoadedAll={isLoadedAll.value}
      loadedTime={loadedTime.value}
      grider={grider.value}
      display={display.value}
    />
  }
})
