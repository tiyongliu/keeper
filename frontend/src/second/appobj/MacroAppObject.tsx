import {defineComponent, inject, PropType, Ref, toRef} from 'vue'
import {filterName} from '/@/second/keeper-tools'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard";
export const extractKey = data => data.name
export const createMatcher = ({ name, title }) => filter => filterName(filter, name, title);

export default defineComponent({
  name: 'MacroAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
  },
  setup(props, {attrs}) {
    const data = toRef(props, 'data')
    const selectedMacro = inject('selectedMacro') as Ref<IPinnedDatabasesItem>
    return () => (
      <AppObjectCore
        {...attrs}
        data={data.value}
        title={data.value!.title}
        icon="img macro"
        onClick={() => selectedMacro.value ? selectedMacro.value = data.value as IPinnedDatabasesItem : null}
      />
    )
  }
})
