import {defineComponent, PropType, ref, toRefs, watch} from 'vue'
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
  emits: ['loadedRows'],
  setup(props, {attrs, emit}) {
    const isLoadedAll = ref<boolean>(false)
    const allRowCount = ref<number | null>(null)
    const errorMessage = ref<string | null>(null)

    const loadedTime = ref(new Date().getTime())

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

    const loadNextDataRef = createRef<boolean>(false)
    const loadedTimeRef = createRef<number | boolean | null>(null)
    const loadedRowsRw= ref(loadedRows.value)
    const handleLoadRowCount = async () => {
      allRowCount.value = await loadRowCount.value!(Object.assign({}, props, attrs))
    }
    const isLoadingRw = ref(isLoading.value)

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

      loadedRowsRw.value = [...loadedRows.value, ...nextRows]
      isLoadedAll.value = nextRows.length === 0

      if (loadNextDataRef.get()) {
        loadNextData()
      }

      emit('loadedRows', loadedRowsRw.value)
    }

    function handleLoadNextData() {
      if (!isLoadedAll.value && !errorMessage.value) {
        if (isFunction(dataPageAvailable.value) && dataPageAvailable.value!(Object.assign({}, props, attrs))) {
          void loadNextData()
        }
      }
    }

    function reload() {
      allRowCount.value = null
      isLoadingRw.value = false
      loadedRowsRw.value = []
      loadedTime.value = new Date().getTime()
      errorMessage.value = null
      loadNextDataRef.set(false)
    }

    watch(() => [display, loadedTime], () => {
      if ((display.value! && display.value?.cache?.refreshTime) > loadedTime.value) {
        reload()
      }
    })

    watch(() => [masterLoadedTime.value, loadedTime.value, display.value], () => {
      if (masterLoadedTime.value && masterLoadedTime.value > loadedTime.value && display.value) {
        display.value.reload()
      }
    })

    return () => (
      <DataGridCore
        {...Object.assign({}, props, attrs)}
        onLoadNextData={handleLoadNextData}
        errorMessage={errorMessage.value}
        isLoading={isLoadingRw.value}
        isLoadedAll={isLoadedAll.value}
        loadedTime={loadedTime.value}
        grider={grider.value}
        display={display.value}
      />
    )
  }
})


