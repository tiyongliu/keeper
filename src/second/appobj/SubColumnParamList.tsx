import {defineComponent, PropType, toRefs, unref,} from 'vue'
import AppObjectList from './AppObjectList'
import {findForeignKeyForColumn} from '/@/packages/tools/src/nameTools'
import ColumnAppObject from '/@/second/appobj/ColumnAppObject'
export default defineComponent({
  name: 'SubColumnParamList',
  props: {
    data: Object as PropType<{ columns: any[] }>,
  },
  setup(props) {
    const {data} = toRefs(props)

    console.log(data, `---------------------`)

    return () => (
      <AppObjectList
        list={(data.value!.columns || []).map(col => ({
          ...data,
          ...col,
          foreignKey: findForeignKeyForColumn(data as any, col)
        }))}
        module={ColumnAppObject}
      />
    )

  }
})
