import {computed, defineComponent, PropType, toRefs} from 'vue'
import {keys} from 'lodash-es'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'
import JsonViewFilterColumn from './JsonViewFilterColumn.vue'
import {GridDisplay} from "/@/second/keeper-datalib";

export default defineComponent({
  name: 'JsonViewFilters',
  components: {
    ManagerInnerContainer
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    useEvalFilters: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  setup(props) {
    const {managerSize, display, isDynamicStructure, useEvalFilters} = toRefs(props)
    const filters = computed<{[uniqueName: string]: string} | null>(() => display.value ? display.value?.config?.filters : null)
    const allFilterNames = computed(() => keys(filters || {}))

    return () => (
      <ManagerInnerContainer width={managerSize.value}>
        {
          allFilterNames.value.map((uniqueName, index) => <JsonViewFilterColumn
            key={index}
            uniqueName={uniqueName}
            display={display.value}
            filters={filters.value!}
            isDynamicStructure={isDynamicStructure.value}
            useEvalFilters={useEvalFilters.value}/>)
        }
      </ManagerInnerContainer>
    )
  }
})
