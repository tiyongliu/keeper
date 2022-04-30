import {defineComponent, PropType} from 'vue'
import DatabaseAppObject from './DatabaseAppObject'
import {IPinnedDatabasesItem} from "/@/second/types/standard.d";

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
    const $props = Object.assign(props, restProps)
    return () => <DatabaseAppObject {...$props} />
  }
})
