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

    <div class="m-1"></div>
    <InlineButton @click="runCommand('new.table')">New table</InlineButton>
    <div class="m-1"></div>
    <InlineButton @click="runCommand('new.collection')">New collection</InlineButton>
  </WidgetsInnerContainer>

  <SearchBoxWrapper v-else>
    <SearchInput placeholder="Search connection or database" v-model:searchValue="filter"/>
    <CloseSearchButton :filter="filter" @close="filter = ''"/>
    <DropDownButton icon="icon plus-thick"/>
    <InlineButton @click="handleRefreshDatabase"
                  title="Refresh database connection and object list">
      <FontIcon icon="icon refresh"/>
    </InlineButton>
  </SearchBoxWrapper>
  <WidgetsInnerContainer>
    <LoadingInfo
      v-if="(status && (status.name == 'pending' || status.name == 'checkStructure' || status.name == 'loadStructure') && objects) || !objects"
      message="Loading database structure" />
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
import {defineComponent, PropType, ref, toRefs, unref, watch} from 'vue';
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
import DatabaseObjectAppObject from '/@/second/appobj/DatabaseObjectAppObject'
import SubColumnParamList from '/@/second/appobj/SubColumnParamList'
import {getObjectTypeFieldLabel} from '/@/second/utility/common'
import {chevronExpandIcon} from '/@/second/icons/expandIcons'
//todo api tables dataSource
import {storeToRefs} from 'pinia'
import {flatten, sortBy} from 'lodash-es'
import {useConnectionInfo, useDatabaseInfo, useDatabaseStatus} from "/@/api/sql"
import {ApplicationDefinition, DatabaseInfo} from '/@/second/keeper-types'
import {filterAppsForDatabase} from '/@/second/utility/appTools'
import {dataBaseStore} from "/@/store/modules/dataBase";

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
    const handleRefreshDatabase = () => {
      // todo apiCall('database-connections/refresh', { conid, database });
    }

    const dataBase = dataBaseStore()
    const {currentDatabase} = storeToRefs(dataBase)
    let objects = ref()
    let status = ref()
    let connection = ref()

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
      useConnectionInfo({conid: unref(conid)}, connection)

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
    }
  }
})
</script>
