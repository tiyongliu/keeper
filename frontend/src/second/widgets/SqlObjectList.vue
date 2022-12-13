<template>
  <WidgetsInnerContainer v-if="status && status.name == 'error'">
    <ErrorInfo :message="`${status.message}`" icon="img error"/>
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>
  </WidgetsInnerContainer>
  <WidgetsInnerContainer v-else-if="objectList.length == 0 &&
  status && status.name != 'pending' && status.name != 'checkStructure' && status.name != 'loadStructure' &&
 objects">
    <ErrorInfo
      :message="`Database ${database} is empty or structure is not loaded, press Refresh button to reload structure`"
      icon="img alert"/>
    <div class="m-1"></div>
    <div class="m-1"></div>
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>
    <template
      v-if="driver && Array.isArray(driver?.databaseEngineTypes) && driver?.databaseEngineTypes?.includes('sql')">
      <div class="m-1"></div>
      <InlineButton @click="runCommand('new.table')">New table</InlineButton>
    </template>
    <template
      v-if="driver && Array.isArray(driver?.databaseEngineTypes) && driver?.databaseEngineTypes?.includes('document')">
      <div class="m-1"></div>
      <InlineButton @click="runCommand('new.collection')">New collection</InlineButton>
    </template>
  </WidgetsInnerContainer>

  <SearchBoxWrapper v-else>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @filter="filter = ''"/>
    <DropDownButton icon="icon plus-thick"/>
    <InlineButton
      @click="handleRefreshDatabase"
      title="Refresh database connection and object list">
      <FontIcon icon="icon refresh"/>
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <LoadingInfo
      v-if="(status && (status.name == 'pending' || status.name == 'checkStructure' || status.name == 'loadStructure') && objects) || !objects"
      message="Loading database structure"/>
    <AppObjectList
      v-else
      :list="objectList.map(x => ({ ...x, conid, database }))"
      :module="databaseObjectAppObject"
      :subItemsComponent="SubColumnParamList"
      :groupFunc="handleGroupFunc"
      :isExpandable="handleExpandable"
      :expandIconFunc="chevronExpandIcon"
      :filter="filter"
      :passProps="{showPinnedInsteadOfUnpin: true}"
    />
  </WidgetsInnerContainer>

</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs, unref, watch, onBeforeUnmount} from 'vue';
import AppObjectList from '/@/second/appobj/AppObjectList'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import SearchBoxWrapper from '/@/second/elements/SearchBoxWrapper.vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton'
import DropDownButton from '/@/second/buttons/DropDownButton'
import runCommand from '/@/second/commands/runCommand'
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import DatabaseObjectAppObject from '/@/second/appobj/DatabaseObjectAppObject'
import SubColumnParamList from '/@/second/appobj/SubColumnParamList'
import {getObjectTypeFieldLabel} from '/@/second/utility/common'
import {chevronExpandIcon} from '/@/second/icons/expandIcons'
import {storeToRefs} from 'pinia'
import {flatten, sortBy} from 'lodash-es'
import {useConnectionInfo, useDatabaseInfo, useDatabaseStatus} from "/@/api/bridge"
import {databaseConnectionsRefreshApi} from '/@/api/simpleApis'
import {ApplicationDefinition, DatabaseInfo} from '/@/second/keeper-types'
import {findEngineDriver} from '/@/second/keeper-tools'
import {filterAppsForDatabase} from '/@/second/utility/appTools'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {useClusterApiStore} from '/@/store/modules/clusterApi'

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
    const {conid, database} = toRefs(props)
    const flag = ref(true)
    const clusterApi = useClusterApiStore()
    const {connection} = storeToRefs(clusterApi)

    const handleRefreshDatabase = async () => {
      try {
        if (flag.value) {
          flag.value = false
          await databaseConnectionsRefreshApi({conid: unref(conid)!, database: unref(database)!})
        }
      } finally {
        flag.value = true
      }
    }

    const bootstrap = useBootstrapStore()
    const {currentDatabase, extensions} = storeToRefs(bootstrap)
    let objects = ref()
    let status = ref()

    const objectList = ref<unknown[]>([])
    const dbApps = ref<ApplicationDefinition[]>([])

    const handleGroupFunc = (data) => {
      return getObjectTypeFieldLabel(unref(data).objectTypeField)
    }

    const handleExpandable = (data) => unref(data).objectTypeField == 'tables' ||
      unref(data).objectTypeField == 'views' || unref(data).objectTypeField == 'matviews'

    watch(() => [conid.value, database.value], () => {
      useDatabaseInfo<DatabaseInfo>({conid: unref(conid), database: unref(database)}, objects)
      useDatabaseStatus<{
        name: 'pending' | 'error' | 'loadStructure' | 'ok';
        counter?: number;
        analysedTime?: number;
      }>({conid: unref(conid), database: unref(database)}, status)
      useConnectionInfo({conid: unref(conid)}, clusterApi.setConnection)

      dbApps.value = filterAppsForDatabase(unref(currentDatabase)?.connection, unref(currentDatabase)!.name, [])
    }, {
      immediate: true
    })

    watch(() => [objects.value, dbApps.value], () => {
      objectList.value = flatten([
        ...['tables', 'collections', 'views', 'matviews', 'procedures', 'functions'].map(objectTypeField =>
          sortBy(
            ((objects.value || {})[objectTypeField] || []).map(obj => ({...obj, objectTypeField})),
            ['schemaName', 'pureName']
          )),
        ...unref(dbApps).map(app => {
          app.queries.map(query => ({
            objectTypeField: 'queries',
            pureName: query.name,
            schemaName: app.name,
            sql: query.sql
          }))
        })
      ])
    })

    const driver = computed(() => extensions.value ? findEngineDriver(connection.value, extensions.value) : null)

    onBeforeUnmount(() => {
      objects.value = null
      status.value = null
      connection.value = null
      objectList.value = []
      dbApps.value = []
    })

    return {
      filter,
      status,
      ...toRefs(props),
      handleRefreshDatabase,
      runCommand,
      objectList,
      objects,
      databaseObjectAppObject: DatabaseObjectAppObject,
      SubColumnParamList,
      handleGroupFunc,
      handleExpandable,
      chevronExpandIcon,
      driver,
    }
  }
})
</script>
