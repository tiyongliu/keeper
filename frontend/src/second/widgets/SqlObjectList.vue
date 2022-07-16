<template>
  <WidgetsInnerContainer v-if="status && status.name == 'error'">
    <ErrorInfo :message="`${status.message}`" icon="img error"/>
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>
  </WidgetsInnerContainer>

  <WidgetsInnerContainer v-else-if="objectList.length == 0 &&
  $status && $status.name != 'pending' && $status.name != 'checkStructure' && $status.name != 'loadStructure' &&
 objects">
    <ErrorInfo
      :message="`Database ${database} is empty or structure is not loaded, press Refresh button to reload structure`"
      icon="img alert"/>
    <div class="m-1"/>
    <div class="m-1"/>
    <InlineButton @click="handleRefreshDatabase">Refresh</InlineButton>

    <div class="m-1"/>
    <InlineButton @click="runCommand('new.table')">New table</InlineButton>

    <div class="m-1"/>
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
    <!--    <LoadingInfo message="Loading database structure" />-->

    <AppObjectList
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
import {computed, defineComponent, PropType, ref, toRef, toRefs, unref, onMounted, watch} from 'vue';
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

import {useDatabaseInfo} from "/@/api/metadataLoaders";
//todo api tables dataSource
import _objectList from './objectList.json'
import _objects from './objects.json'


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
    const _conid = toRef(props, 'conid')
    const _database = toRef(props, 'database')
    //todo
    // $: status = useDatabaseStatus({ conid, database });

    // const status = computed(() => {})

    const handleRefreshDatabase = () => {
      // todo apiCall('database-connections/refresh', { conid, database });
    }

    const objectList = computed<any[]>(() => _objectList)
    const objects = computed<{ tables: any[] }>(() => _objects)

    const handleGroupFunc = (data) => {
      return getObjectTypeFieldLabel(unref(data).objectTypeField)
    }

    const handleExpandable = (data) => unref(data).objectTypeField == 'tables' ||
      unref(data).objectTypeField == 'views' || unref(data).objectTypeField == 'matviews'


    //useDatabaseInfo

    const showDatabaseInfo = async ({conid, database}) => {
      const data = await useDatabaseInfo({conid, database})
      console.log(data, `dddddddddddddddddddddddddddddddddddddddd`)
    }

    onMounted(() => {
      void showDatabaseInfo({conid: _conid.value, database: _database.value})
    })

    watch(() => [_conid.value, _database.value], (ar) => {
      console.log(ar, `111111111111111111111111111111111`)
      const [c, d] = ar
      void showDatabaseInfo({conid: c, database: d})
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
      chevronExpandIcon
    }
  }
})
</script>
