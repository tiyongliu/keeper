<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:value="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <InlineButton title="Add new connection">
      <plus-outlined/>
    </InlineButton>

    <InlineButton title="Add new connection">
      <redo-outlined/>
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <AppObjectList />
  </WidgetsInnerContainer>
</template>

<script lang="ts">
  import {defineComponent, ref} from 'vue';
  import SearchBoxWrapper from './SearchBoxWrapper.vue'
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  import SearchInput from '/@/second/elements/SearchInput.vue'
  import InlineButton from '/@/second/buttons/InlineButton.vue'
  import CloseSearchButton from '/@/second/buttons/CloseSearchButton.vue'
  import AppObjectList from '/@/second/appobj/AppObjectList.vue'
  import {PlusOutlined, RedoOutlined} from '@ant-design/icons-vue'

  export default defineComponent({
    name: "ConnectionList",
    components: {
      SearchBoxWrapper,
      WidgetsInnerContainer,
      CloseSearchButton,
      SearchInput,
      InlineButton,
      RedoOutlined,
      PlusOutlined,
      AppObjectList,
    },
    setup() {
      const hidden = ref(false)
      const filter = ref('')

      return {
        hidden,
        filter
      }
    }
  })
</script>

<style scoped>

</style>
