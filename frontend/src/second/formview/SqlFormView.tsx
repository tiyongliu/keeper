import {computed, defineComponent, PropType, ref, toRefs, watch, watchEffect} from 'vue'
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
    referenceSourceChanged: {
      type: Function as PropType<(value: any[], loadedTime: number) => void>
    }
  },
  setup(props, {attrs}) {
    const {
      changeSetState,
      dispatchChangeSet,
      formDisplay,
      conid,
      database,
      masterLoadedTime,
      referenceSourceChanged
    } = toRefs(props)
    const isLoadingData = ref(false)
    const isLoadedData = ref(false)
    const rowData = ref<Nullable<any>>(null)
    const isLoadingCount = ref(false)
    const isLoadedCount = ref(false)
    const loadedTime = ref(new Date().getTime())
    const allRowCount = ref<Nullable<number>>(null)
    const rowCountBefore = ref<Nullable<number>>(null)


    const former = computed(() => (rowData.value && formDisplay.value)
      ? new ChangeSetFormer(rowData.value, changeSetState.value, dispatchChangeSet.value, formDisplay.value!)
      : null)

    const handleLoadCurrentRow = async () => {
      if (isLoadingData.value) return;
      let newLoadedRow = false;
      if (formDisplay.value && (formDisplay.value.config.formViewKeyRequested || formDisplay.value.config.formViewKey)) {
        isLoadingData.value = true
        const row = await loadRow(Object.assign({}, props, attrs), formDisplay.value.getCurrentRowQuery())
        isLoadingData.value = false;
        isLoadedData.value = true;
        rowData.value = row;
        loadedTime.value = new Date().getTime();
        newLoadedRow = row;
      }
      if (formDisplay.value && formDisplay.value.config.formViewKeyRequested && newLoadedRow) {
        formDisplay.value.cancelRequestKey(newLoadedRow);
      }
      if (!newLoadedRow && !(formDisplay.value && formDisplay.value.config.formViewKeyRequested)) {
        await handleNavigate('first');
      }
    }

    const handleNavigate = async command => {
      isLoadingData.value = true
      const row = await loadRow(Object.assign({}, props, attrs), formDisplay.value!.navigateRowQuery(command))
      if (row) {
        formDisplay.value!.navigate(row)
      }
      isLoadingData.value = false
      isLoadedData.value = true
      isLoadingCount.value = false
      isLoadedCount.value = false
      allRowCount.value = null
      rowCountBefore.value = null
      rowData.value = row

      loadedTime.value = new Date().getTime()
    }

    const handleLoadRowCount = async () => {
      isLoadingCount.value = true;
      const countRow = await loadRow(Object.assign({}, props, attrs), formDisplay.value!.getCountQuery());
      const countBeforeRow = await loadRow(Object.assign({}, props, attrs), formDisplay.value!.getBeforeCountQuery());

      isLoadedCount.value = true;
      isLoadingCount.value = false;
      allRowCount.value = countRow ? parseInt(countRow.count) : null;
      rowCountBefore.value = countBeforeRow ? parseInt(countBeforeRow.count) : null;
    }

    watch(() => [masterLoadedTime.value, loadedTime.value], () => {
      if (masterLoadedTime.value && masterLoadedTime.value > loadedTime.value) {
        formDisplay.value!.reload()
      }
    })

    watchEffect(() => {
      if (formDisplay.value && formDisplay.value.isLoadedCorrectly) {
        if (!isLoadedData.value && !isLoadingData.value) void handleLoadCurrentRow()
        if (isLoadedData.value && !isLoadingCount.value && !isLoadedCount.value) void handleLoadRowCount()
      }

      if (referenceSourceChanged.value && rowData.value) {
        referenceSourceChanged.value([rowData.value], loadedTime.value)
      }
    })

    return () => (
      <FormView
        {...Object.assign({
          formDisplay: formDisplay.value,
          conid: conid.value,
          database: database.value,
        }, attrs)}
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
  if (response.errorMessage) return response
  return response.rows[0]
}
