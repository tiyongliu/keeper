import {computed, defineComponent, PropType, toRefs, unref,} from 'vue'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {TableInfo} from "/@/packages/tools/types/dbinfo";
import {_getColumnIcon} from '/@/second/elements/ColumnLabel_'

export default defineComponent({
  name: 'ColumnAppObject',
  props: {
    data: Object as PropType<{ foreignKey: TableInfo, columnName: string, refTableName?: string, dataType?: string }>,
  },
  setup(props, {attrs}) {
    const {data} = toRefs(props)

    const extInfo = computed(() => unref(data)!.foreignKey ? `${unref(data)!.dataType} -> ${unref(data)!.refTableName}` : unref(data)!.dataType)

    return () => (
      <AppObjectCore
        {...attrs}
        data={unref(data)}
        title={unref(data)!.columnName}
        extInfo={unref(extInfo)}
        icon={_getColumnIcon(unref(data), true)}
      />
    )
  }
})
