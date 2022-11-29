import {computed, defineComponent, PropType, ref, toRefs, watch} from 'vue'
import {databaseConnectionsSqlSelectApi} from '/@/api/simpleApis'
import FormView from './FormView.vue'
import ChangeSetFormer from './ChangeSetFormer'
import {TableFormViewDisplay} from '/@/second/keeper-datalib'

export default defineComponent({
  name: 'SqlFormView',
  props: {
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
    masterLoadedTime: {
      type: Number as PropType<number>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
  },
  setup(props, {attrs}) {
    const {changeSetState, dispatchChangeSet, formDisplay, masterLoadedTime} = toRefs(props)
    const isLoadingData = ref(false)
    const isLoadedData = ref(false)
    const rowData = ref<Nullable<any>>(null)
    const isLoadedCount = ref(false)
    const loadedTime = ref(new Date().getTime())
    const allRowCount = ref<Nullable<number>>(null)
    const rowCountBefore = ref<Nullable<number>>(null)

    watch(() => [masterLoadedTime.value, loadedTime.value], () => {
      if (masterLoadedTime.value && masterLoadedTime.value > loadedTime.value) {
        formDisplay.value!.reload()
      }
    })

    const former = computed(() => (rowData.value && formDisplay.value)
      ? new ChangeSetFormer(rowData.value, changeSetState.value, dispatchChangeSet.value, formDisplay.value!)
      : null)

    const handleNavigate = async command => {
      isLoadingData.value = true
      const row = await loadRow(Object.assign(props, attrs), formDisplay.value!.navigateRowQuery(command))
      if (row) {
        formDisplay.value!.navigate(row)
      }
      isLoadingData.value = false
      isLoadedData.value = true
      isLoadedCount.value = false
      allRowCount.value = null
      rowCountBefore.value = null
      rowData.value = row
      loadedTime.value = new Date().getTime()
    }

    return () => (
      <FormView
        {...Object.assign({}, props, attrs)}
        former={former.value}
        isLoading={isLoadingData.value}
        allRowCount={allRowCount.value}
        rowCountBefore={rowCountBefore.value}
        navigate={handleNavigate}
      />
    )
  }
})

async function loadRow(props, select) {
  const {conid, database} = props;

  if (!select) return null;

  const response = await databaseConnectionsSqlSelectApi({
    conid,
    database,
    select,
  }) as any
  console.log(response, `rrrrrrrrrrrrrrrrrrrrrr`)
  if (response.errorMessage) return response
  return response.rows[0]
}
