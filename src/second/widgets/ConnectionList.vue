<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <InlineButton title="Add new connection">
      <FontIcon icon="icon plus-thick" />
    </InlineButton>

    <InlineButton title="Add new connection">
      <FontIcon icon="icon refresh" />
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <AppObjectList
      :list="sortBy(connectionsWithStatus, connection => (getConnectionLabel(connection) || '').toUpperCase())"
      :filter="filter"
      :expandOnClick="false"
      :isExpandable="(data) => [
        'b0cf1450-a66d-11ec-a868-3720e8369945',
        '065caa90-a8c6-11ec-9b4b-6f98950c4d7a'
      ].includes(data._id) && !data.singleDatabase"
    />
  </WidgetsInnerContainer>
</template>

<script lang="ts">
  import {defineComponent, ref, onMounted} from 'vue';
  import {sortBy} from 'lodash-es'
  import SearchBoxWrapper from './SearchBoxWrapper.vue'
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  import SearchInput from '/@/second/elements/SearchInput.vue'
  import InlineButton from '/@/second/buttons/InlineButton.vue'
  import CloseSearchButton from '/@/second/buttons/CloseSearchButton.vue'
  // import AppObjectList from '/@/second/appobj/AppObjectList.vue'
  import AppObjectList from '/@/second/appobj/AppObjectList'

  import FontIcon from '/@/second/icons/FontIcon.vue'
  import getConnectionLabel from '/@/second/utility/getConnectionLabel'
  export default defineComponent({
    name: "ConnectionList",
    components: {
      SearchBoxWrapper,
      WidgetsInnerContainer,
      CloseSearchButton,
      SearchInput,
      InlineButton,
      AppObjectList,
      FontIcon,
    },
    setup() {
      const hidden = ref(false)
      const filter = ref('')
      const connectionsWithStatus = [{"server":"localhost","engine":"mongo@dbgate-plugin-mongo","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","useDatabaseUrl":"","_id":"b0cf1450-a66d-11ec-a868-3720e8369945","status":{"name":"ok"}}, {"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a"}]

      onMounted(() => {

      })

      return {
        hidden,
        filter,
        connectionsWithStatus,
        sortBy,
        getConnectionLabel,

      }
    }
  })
</script>

<style scoped>

</style>
