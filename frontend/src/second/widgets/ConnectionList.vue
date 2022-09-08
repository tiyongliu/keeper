<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <InlineButton title="Add new connection" @click="openModal">
      <FontIcon icon="icon plus-thick"/>
    </InlineButton>
    <InlineButton title="Add new connection" @click="handleRefreshConnections">
      <FontIcon icon="icon refresh"/>
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <AppObjectList
      v-if="Array.isArray(connectionsWithStatus) && connectionsWithStatus.length > 0"
      :list="sortBy(connectionsWithStatus, connection => (getConnectionLabel(connection) || '').toUpperCase())"
      :filter="filter"
      :module="connectionAppObject"
      :subItemsComponent="SubDatabaseList"
      expandOnClick
      :isExpandable="handleExpandable"
      :passProps="{showPinnedInsteadOfUnpin: true}"
    />
    <LargeButton
      v-else
      icon="icon new-connection"
      fillHorizontal
      @visible="openModal">
      Add new connection
    </LargeButton>
    <ConnectionModal @register="register" @closeCurrentModal="closeModal"/>
  </WidgetsInnerContainer>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref, unref, watch} from 'vue'
import {storeToRefs} from 'pinia'
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
import LargeButton from '/@/second/buttons/LargeButton.vue'

//TODO
import ConnectionModal from '/@/second/modals/ConnectionModal.vue'
import {useModal} from '/@/components/Modal'
import {handleRefreshApi} from '/@/api/serverConnections'
import {useConnectionList, useServerStatus} from '/@/api/sql'
import {IActiveConnection, IConnectionStatus} from '/@/second/typings/types/connections.d'

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
    LargeButton,
    ConnectionModal,
  },
  setup() {
    const hidden = ref(false)
    const filter = ref('')
    const dataBase = dataBaseStore()
    const {openedConnections} = storeToRefs(dataBase)

    const handleExpandable = (data) => unref(openedConnections).includes(unref(data)._id)
      && !unref(data).singleDatabase
    const connectionsWithStatus = ref<IActiveConnection[]>([])

    const connections = ref()
    const serverStatus = ref()

    onMounted(() => {
      useConnectionList<IActiveConnection[]>(connections)
      useServerStatus<{ [key: string]: IConnectionStatus }>(serverStatus)
    })

    watch(() => [connections, serverStatus], () => {
      connectionsWithStatus.value =
        connections.value && serverStatus.value ?
          connections.value.map(conn => ({...conn, status: serverStatus.value[conn._id]})) :
          connections.value
    }, {
      deep: true,
      immediate: true,
    })

    const handleRefreshConnections = () => {
      for (const conid of unref(openedConnections)) {
        void handleRefreshApi({conid})
      }
    }

    const [register, {openModal, closeModal}] = useModal()
    return {
      hidden,
      filter,
      connectionsWithStatus,
      sortBy,
      getConnectionLabel,
      connectionAppObject: ConnectionAppObject,
      SubDatabaseList,
      handleExpandable,
      runCommand,
      register,
      openModal,
      closeModal,
      connections, serverStatus,
      handleRefreshConnections
    }
  }
})
</script>
