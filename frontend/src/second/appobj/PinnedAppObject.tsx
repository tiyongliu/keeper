import {defineComponent, PropType} from 'vue'
import DatabaseAppObject from './DatabaseAppObject'
import DatabaseObjectAppObject from './DatabaseObjectAppObject'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d";

export const extractKey = props => props.name

export default defineComponent({
  name: 'PinnedAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
  },
  setup(props, {attrs}) {
    const {onClick, expandIcon, onExpand, ...restProps} = attrs
    const $props = Object.assign(props, restProps) as any

    return () => $props.data.objectTypeField ?
      <DatabaseObjectAppObject {...$props}/> :
      <DatabaseAppObject {...$props} />
  }
})
