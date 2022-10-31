import {defineComponent, PropType, ref, toRefs, unref, watch, onMounted} from 'vue'
import {getIntSettingsValue} from '/@/second/settings/settingsTools'
import createRef from '/@/second/utility/createRef'
import DataGridCore from './DataGridCore.vue'
import {GridDisplay} from "/@/second/keeper-datalib";

export default defineComponent({
  name: 'LoadingDataGridCore',
  props: {
    loadDataPage: {
      type: Function as PropType<(props: any, offset: any, limit: any) => Promise<any[]>>,
    },
    isLoading: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: []
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
  },
  setup(props, {attrs}) {
    const allRowCount = ref(null)
    const loadedTime = ref(new Date().getTime())

    const {isLoading, loadDataPage, loadedRows, display} = toRefs(props)

    const loadedTimeRef = createRef<Number | null>(null)

    const handleLoadRowCount = async () => {

    }

    async function loadNextData() {
      if (isLoading.value) return
      isLoading.value = true

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

      isLoading.value = false
    }

    function handleLoadNextData() {
      void loadNextData()
      console.log(``)
    }

    function reload() {
      allRowCount.value = null
      isLoading.value = false
      loadedRows.value = []
      loadedTime.value = new Date().getTime()
    }

    watch(() => [display, loadedTime], () => {
      // @ts-ignore
      if (display.value?.cache?.refreshTime > loadedTime.value) {
        reload()
      }
    })

    watch(() => display.value, () => {
      console.log(display.value, `displaydisplaydisplaydisplaydisplaydisplay`)
    })

    return () => (
      <DataGridCore display={display.value} onLoadNextData={handleLoadNextData} />
    )
  }
})


