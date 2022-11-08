<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <InlineButton @click="showModal">Add</InlineButton>
    <InlineButton>Hide</InlineButton>
    <InlineButton>Show</InlineButton>
  </SearchBoxWrapper>
  <ManagerInnerContainer :width="managerSize">
    <ColumnManagerRow />
  </ManagerInnerContainer>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRef} from 'vue'
import SearchBoxWrapper from '/@/second/elements/SearchBoxWrapper.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'
import ColumnManagerRow from '/@/second/datagrid/ColumnManagerRow.vue'

export default defineComponent({
  name: "ColumnManager",
  components: {
    SearchBoxWrapper,
    SearchInput,
    CloseSearchButton,
    InlineButton,
    ManagerInnerContainer,
    ColumnManagerRow
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    }
  },
  setup(props) {
    const filter = ref('')
    const selectedColumns = ref<unknown[]>([])
    const currentColumnUniqueName = ref()

    function showModal() {}

    const managerSize = toRef(props, 'managerSize')

    function setSelectedColumns(value: unknown[]) {
      selectedColumns.value = value
      if (value.length > 0) {
        currentColumnUniqueName.value = value[0]
      }
    }

    return {
      filter,
      managerSize,
      showModal,
      setSelectedColumns
    }
  }
})
</script>

<style scoped>
.focus-field {
  position: absolute;
  left: -1000px;
  top: -1000px;
}
</style>
