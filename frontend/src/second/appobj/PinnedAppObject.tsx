import {defineComponent, PropType, toRefs} from 'vue'
import DatabaseAppObject from './DatabaseAppObject'
import DatabaseObjectAppObject from './DatabaseObjectAppObject'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d";
import {omit} from 'lodash-es'
export const extractKey = props => props.name

export default defineComponent({
  name: 'PinnedAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
  },
  setup(props, {attrs}) {
    const {data} = toRefs(props)
    const restProps = omit(attrs, ['onClick', 'expandIcon', 'onExpand'])
    const $props = Object.assign({}, props, restProps) as any

    // 使用下面这种，会解决控制板警告，但是取消表收藏会有bug，抽空需要学习vue源码
    // const $props = shallowReadonly({
    //   ...props, ...restProps
    // }) as any

    return () => data.value && (
      <>
        {data.value?.objectTypeField
          ? <DatabaseObjectAppObject {...$props}/>
          : <DatabaseAppObject {...$props} />}
      </>
    )
  }
})
