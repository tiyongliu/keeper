<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:value="filter"/>
    <CloseSearchButton :filter="filter" @filter="filter = ''"/>
    <InlineButton @click="handleAddKey" title="Add new key">
      <FontIcon icon="icon plus-thick"/>
    </InlineButton>
    <InlineButton @click="handleRefreshDatabase" title="Refresh key list">
      <FontIcon icon="icon refresh" />
    </InlineButton>
    <WidgetsInnerContainer>
      DbKeysSubTree component
    </WidgetsInnerContainer>
  </SearchBoxWrapper>
</template>

<script lang="ts">
import {ref} from 'vue'
import SearchBoxWrapper from '/@/second/widgets/SearchBoxWrapper.vue'
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
export default {
  name: 'DbKeysTree',
  components: {
    SearchBoxWrapper,
    WidgetsInnerContainer,
    SearchInput,
    CloseSearchButton,
    InlineButton,
    FontIcon,
  },
  setup() {
    const filter = ref('')

    function handleAddKey() {}

    function handleRefreshDatabase() {}
    return {
      filter,
      handleAddKey,
      handleRefreshDatabase
    }
  }
}
</script>
