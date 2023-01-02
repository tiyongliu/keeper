import {defineComponent, PropType, ref, toRefs, watchEffect} from 'vue'
//uniqWith, isEqual
import {getIntSettingsValue} from '/@/second/settings/settingsTools'
import createRef from '/@/second/utility/createRef'
import DataGridCore from './DataGridCore.vue'
import {GridDisplay} from '/@/second/keeper-datalib'
import {isFunction} from '/@/utils/is'
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
    isLoading: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: []
    },
    grider: {
      type: Object as PropType<Grider>
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    masterLoadedTime: {
      type: Number as PropType<number>,
    }
  },
  emits: ['selectedCellsPublished', 'update:loadedRows'],
  setup(props, {attrs, emit}) {
    const {
      isLoading,
      loadDataPage,
      loadedRows,
      loadRowCount,
      display,
      dataPageAvailable,
      grider,
      masterLoadedTime
    } = toRefs(props)

    const isLoadedAll = ref<boolean>(false)
    const allRowCount = ref<Nullable<number>>(null)
    const errorMessage = ref<Nullable<string>>(null)
    const loadedTime = ref(new Date().getTime())
    const domGrid = ref<Nullable<HTMLElement>>(null)
    const loadNextDataRef = createRef<boolean>(false)
    const loadedTimeRef = createRef<number | boolean | null>(null)
    const loadedRowsRw = ref(loadedRows.value)
    const isLoadingRw = ref(isLoading.value)

    const handleLoadRowCount = async () => {
      allRowCount.value = await loadRowCount.value!(Object.assign({}, props, attrs))
    }

    async function loadNextData() {
      if (isLoadingRw.value) return
      loadedTimeRef.set(false)
      isLoadingRw.value = true

      const loadStart = new Date().getTime()

      loadedTimeRef.set(loadStart)

      const nextRows = await loadDataPage.value!(
        Object.assign({}, props, attrs),
        loadedRows.value.length,
        getIntSettingsValue('dataGrid.pageSize', 100, 5, 1000)
      )

      if (loadedTimeRef.get() !== loadStart) {
        return
      }

      isLoadingRw.value = false

      if (nextRows.errorMessage) {
        errorMessage.value = nextRows.errorMessage
      } else {
        if (allRowCount.value == null) await handleLoadRowCount()
      }
      console.log(`loadedRows`, loadedRows.value)
      console.log(`nextRows`, nextRows)
      loadedRowsRw.value = [...loadedRows.value, ...nextRows]
      isLoadedAll.value = nextRows.length === 0
      if (loadNextDataRef.get()) {
         loadNextData()
      }
      emit('update:loadedRows', loadedRowsRw.value)
    }

    function handleLoadNextData() {
      if (!isLoadedAll.value && !errorMessage.value && (grider.value && !grider.value.disableLoadNextPage)) {
        if (isFunction(dataPageAvailable.value) && dataPageAvailable.value!(Object.assign({}, props, attrs))) {
          void loadNextData()
        }
      }
    }

    function selectedCellsPublished(data) {
      emit('selectedCellsPublished', data)
    }

    function reload() {
      allRowCount.value = null
      isLoadingRw.value = false
      loadedRowsRw.value = []
      loadedTime.value = new Date().getTime()
      errorMessage.value = null
      loadNextDataRef.set(false)
    }

    watchEffect(() => {
      if ((display.value! && display.value?.cache?.refreshTime) > loadedTime.value) {
        reload()
      }

      if (masterLoadedTime.value && masterLoadedTime.value > loadedTime.value && display.value) {
        display.value.reload()
      }
    })

    return () => <DataGridCore
      {...Object.assign({}, props, attrs)}
      ref={domGrid}
      onSelectedCellsPublished={selectedCellsPublished}
      onLoadNextData={handleLoadNextData}
      errorMessage={errorMessage.value}
      isLoading={isLoadingRw.value}
      isLoadedAll={isLoadedAll.value}
      loadedTime={loadedTime.value}
      grider={grider.value}
      display={display.value}
    />
  }
})
