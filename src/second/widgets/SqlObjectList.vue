<template>
  <!--<WidgetsInnerContainer v-if="status && status.name == 'error'">
    <ErrorInfo :message="`${status.message}`" icon="img error"/>
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>
  </WidgetsInnerContainer>

  <WidgetsInnerContainer>
    <ErrorInfo
      :message="`Database ${database} is empty or structure is not loaded, press Refresh button to reload structure`"
      icon="img alert"/>
    <div class="m-1" />
    <div class="m-1" />
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>

    <div class="m-1" />
    <InlineButton @click="runCommand('new.table')">New table</InlineButton>

    <div class="m-1" />
    <InlineButton @click="runCommand('new.collection')">New collection</InlineButton>
  </WidgetsInnerContainer>
-->

  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <DropDownButton icon="icon plus-thick"/>
    <InlineButton @click="handleRefreshDatabase" title="Refresh database connection and object list">
      <FontIcon icon="icon refresh" />
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <LoadingInfo message="Loading database structure" />
  </WidgetsInnerContainer>
</template>

<script lang="ts">
  import {defineComponent, PropType, computed, watch, unref, toRefs, ref} from 'vue';
  import AppObjectList from '/@/second/appobj/AppObjectList'
  import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
  import FontIcon from '/@/second/icons/FontIcon.vue'
  import InlineButton from '/@/second/buttons/InlineButton.vue'
  import SearchBoxWrapper from '/@/second/elements/SearchBoxWrapper.vue'
  import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
  import SearchInput from '/@/second/elements/SearchInput.vue'
  import CloseSearchButton from '/@/second/buttons/CloseSearchButton.vue'
  import DropDownButton from '/@/second/buttons/DropDownButton'
  import runCommand from '/@/second/commands/runCommand'
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  export default defineComponent({
    name: "SqlObjectList",
    props: {
      conid: {
        type: String as PropType<string>
      },
      database: {
        type: String as PropType<string>
      }
    },
    components: {
      AppObjectList,
      WidgetsInnerContainer,
      LoadingInfo,
      ErrorInfo,
      FontIcon,
      InlineButton,
      SearchBoxWrapper,
      SearchInput,
      CloseSearchButton,
      DropDownButton,
    },
    setup(props) {
      const filter = ref('')
      //todo
      // $: status = useDatabaseStatus({ conid, database });

      const status = computed(() => {})

      const handleRefreshDatabase = () => {
        // todo apiCall('database-connections/refresh', { conid, database });
      }

      return {
        filter,
        status,
        ...toRefs(props),
        handleRefreshDatabase,
        runCommand,
      }
    }
  })
</script>
