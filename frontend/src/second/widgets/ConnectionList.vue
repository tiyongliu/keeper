<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <InlineButton title="Add new connection" @click="runCommand('new.connection')">
      <FontIcon icon="icon plus-thick"/>
    </InlineButton>

    <InlineButton title="Add new connection">
      <FontIcon icon="icon refresh"/>
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <AppObjectList
      :list="sortBy(connectionsWithStatus, connection => (getConnectionLabel(connection) || '').toUpperCase())"
      :filter="filter"
      :module="connectionAppObject"
      :subItemsComponent="SubDatabaseList"
      expandOnClick
      :isExpandable="handleExpandable"
      :passProps="{showPinnedInsteadOfUnpin: true}"
    />
  </WidgetsInnerContainer>
</template>

<script lang="ts">
import {defineComponent, ref, unref} from 'vue'
import {sortBy} from 'lodash-es'
import SearchBoxWrapper from '/@/second/widgets/SearchBoxWrapper.vue'
import WidgetsInnerContainer from '/@/second/widgets//WidgetsInnerContainer.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import AppObjectList from '/@/second/appobj/AppObjectList'
import FontIcon from '/@/second/icons/FontIcon.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'
import ConnectionAppObject from '/@/second/appobj/ConnectionAppObject'
import SubDatabaseList from '/@/second/appobj/SubDatabaseList'
import {dataBaseStore} from '/@/store/modules/dataBase'
import runCommand from '/@/second/commands/runCommand'

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
    const dataBase = dataBaseStore()

    const connectionsWithStatus = [{
      "server": "localhost",
      "engine": "mysql@dbgate-plugin-mysql",
      "sshMode": "userPassword",
      "sshPort": "22",
      "sshKeyfile": "C:\\Users\\Administrator\\.ssh\\id_rsa",
      "user": "root",
      "password": "crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==",
      "_id": "065caa90-a8c6-11ec-9b4b-6f98950c4d7a",
      "status": {"name": "ok"}
    }]

    const handleExpandable = (data) => dataBase.$state.openedConnections.includes(unref(data)._id)
      && !unref(data).singleDatabase

    return {
      hidden,
      filter,
      connectionsWithStatus,
      sortBy,
      getConnectionLabel,
      connectionAppObject: ConnectionAppObject,
      SubDatabaseList,
      handleExpandable,
      runCommand
    }
  }
})
</script>