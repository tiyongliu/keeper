import {defineComponent, inject, PropType, Ref, toRef} from 'vue'
import {filterName} from '/@/second/keeper-tools'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {Nullable} from '/@/utils/types'
import {MacroDefinition} from '/@/second/keeper-datalib'

export const extractKey = data => data.name
export const createMatcher = ({ name, title }) => filter => filterName(filter, name, title);

export default defineComponent({
  name: 'MacroAppObject',
  props: {
    data: {
      type: Object as PropType<MacroDefinition>,
    },
  },
  setup(props, {attrs}) {
    const data = toRef(props, 'data')
    const selectedMacro = inject('selectedMacro') as Ref<Nullable<MacroDefinition>>

    function changeSelectedMacro() {
      if (selectedMacro.value && data.value) {
        selectedMacro.value = data.value
      }
    }

    return () => (
      <AppObjectCore
        {...attrs}
        data={data.value}
        title={data.value!.title}
        icon="img macro"
        onClick={changeSelectedMacro}
      />
    )
  }
})
