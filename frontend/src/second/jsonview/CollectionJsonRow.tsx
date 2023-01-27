import {computed, defineComponent, PropType, toRefs} from 'vue'
import {isNumber} from 'lodash-es'
import JSONTree from '/@/second/jsontree/JSONTree.vue'
import Grider from '/@/second/datagrid/Grider'

export default defineComponent({
  name: 'CollectionJsonRow',
  props: {
    rowIndex: {
      type: Number as PropType<number>,
    },
    grider: {
      type: Object as PropType<Grider>,
    },
  },
  setup(props) {
    const {rowIndex, grider} = toRefs(props)
    const rowData = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowData(rowIndex.value)
      }
      return null
    })

    const rowStatus = computed(() => {
      if (grider.value && isNumber(rowIndex.value)) {
        return grider.value.getRowStatus(rowIndex.value)
      }
      return null
    })

    return () => (
      <JSONTree
        value={rowData.value}
        labelOverride={`(${rowIndex.value! + 1})`}
        isModified={rowStatus.value && rowStatus.value.status == 'updated'}
        isInserted={rowStatus.value && rowStatus.value.status == 'inserted'}
        isDeleted={rowStatus.value && rowStatus.value.status == 'deleted'}
      />
    )
  }
})
