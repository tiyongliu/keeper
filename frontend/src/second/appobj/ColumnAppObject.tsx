import {computed, defineComponent, PropType, toRef} from 'vue'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {getColumnIcon} from '/@/second/elements/ColumnLabel.vue'

export default defineComponent({
  name: 'ColumnAppObject',
  props: {
    data: Object as PropType<{ foreignKey: { refTableName: string }, columnName: string, refTableName?: string, dataType?: string }>,
  },
  setup(props, {attrs}) {
    const data = toRef(props, 'data')

    const extInfo = computed(() => data.value!.foreignKey ? `${data.value!.dataType} -> ${data.value!.foreignKey.refTableName}` : data.value!.dataType)

    return () => (
      <AppObjectCore
        {...attrs}
        data={data.value}
        title={data.value!.columnName}
        extInfo={extInfo.value}
        icon={getColumnIcon(data.value, true)}
      />
    )
  }
})
