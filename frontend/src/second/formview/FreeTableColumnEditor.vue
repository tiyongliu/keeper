<template>
  <ManagerInnerContainer :width="managerSize">
    <template v-for="(column, index) in structure?.columns || []">
      <ColumnNameEditor
        v-if="index == editingColumn"
        :defaultValue="column.columnName"
        focusOnCreate
        blurOnEnter
      />
      <ColumnManagerRow
        v-else
        :column="column"
      />
    </template>
    <ColumnNameEditor
      placeholder="New column"
      :existingNames="(structure?.columns || []).map(x => x.columnName)" />
  </ManagerInnerContainer>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs, computed, ref} from 'vue'
import {TableFormViewDisplay} from '/@/second/keeper-datalib'
import {DatabaseInfo} from '/@/second/keeper-types'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'
import ColumnNameEditor from '/@/second/freetable/ColumnNameEditor.vue'
import ColumnManagerRow from '/@/second/freetable/ColumnManagerRow.vue'
export default defineComponent({
  name: 'FreeTableColumnEditor',
  components: {
    ManagerInnerContainer,
    ColumnNameEditor,
    ColumnManagerRow
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    setConfig: {
      type: Function as PropType<(target: any) => void>
    },

    modelState: {
      type: Object as PropType<{structure: DatabaseInfo & {columns: any[]}}>
    },
  },
  setup(props){
    const {managerSize,modelState} = toRefs(props)
    const editingColumn = ref<Nullable<HTMLElement>>(null)
    const structure = computed(() => modelState.value ? modelState.value?.structure : null)

    return {
      managerSize,
      structure,
      editingColumn
    }
  }
})
</script>
