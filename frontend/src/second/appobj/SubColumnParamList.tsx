import {defineComponent, PropType, toRef} from 'vue'
import AppObjectList from './AppObjectList'
import {findForeignKeyForColumn} from 'keeper-tools'
import ColumnAppObject from '/@/second/appobj/ColumnAppObject'
export default defineComponent({
  name: 'SubColumnParamList',
  props: {
    data: Object as PropType<{ columns: any[] }>,
  },
  setup(props) {
    const data = toRef(props, 'data')

    return () => (
      <AppObjectList
        list={(data.value!.columns || []).map(col => ({
          ...data.value,
          ...col,
          foreignKey: findForeignKeyForColumn(data.value as any, col)
        }))}
        module={ColumnAppObject}
      />
    )

  }
})
