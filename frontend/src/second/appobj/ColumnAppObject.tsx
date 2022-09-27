import {computed, defineComponent, PropType, toRefs, unref,} from 'vue'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {getColumnIcon} from '/@/second/elements/ColumnLabel.vue'

export default defineComponent({
  name: 'ColumnAppObject',
  props: {
    data: Object as PropType<{ foreignKey: {refTableName: string}, columnName: string, refTableName?: string, dataType?: string }>,
  },
  setup(props, {attrs}) {
    const {data} = toRefs(props)

    const extInfo = computed(() => unref(data)!.foreignKey ? `${unref(data)!.dataType} -> ${unref(data)!.foreignKey.refTableName}` : unref(data)!.dataType)

    return () => (
      <AppObjectCore
        {...attrs}
        data={unref(data)}
        title={unref(data)!.columnName}
        extInfo={unref(extInfo)}
        icon={getColumnIcon(unref(data), true)}
      />
    )
  }
})
