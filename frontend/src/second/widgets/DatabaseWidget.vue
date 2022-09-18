<template>
  <WidgetColumnBar :hidden="hidden">
    <!--这个是上面数据库 及 db列表-->
    <WidgetColumnBarItem title="Connections" name="connections" height="35%" storageName="connectionsWidget">
      <ConnectionList/>
    </WidgetColumnBarItem>
    <WidgetColumnBarItem
      title="Pinned"
      name="pinned"
      height="15%"
      storageName="pinnedItemsWidget"
      :skip="!pinnedDatabases?.length &&
        !pinnedTables.some(x => x.conid == currentDatabase.connection._id && x.database == currentDatabase?.name)"
    >
      <PinnedObjectsList/>
    </WidgetColumnBarItem>
    <!--数据库 table 列表-->
    <WidgetColumnBarItem
      v-if="conid && (database || singleDatabase)"
      :title="driver && Array.isArray(driver?.databaseEngineTypes) && driver?.databaseEngineTypes?.includes('document')
      ? 'Collections' : 'Tables, views, functions'"
      name="dbObjects"
      storageName="dbObjectsWidget">
      <SqlObjectList :conid="conid" :database="database"/>
    </WidgetColumnBarItem>
    <WidgetColumnBarItem v-else title="Database content" name="dbObjects" storageName="dbObjectsWidget">
      <WidgetsInnerContainer>
        <ErrorInfo message="Database not selected" icon="img alert"/>
      </WidgetsInnerContainer>
    </WidgetColumnBarItem>
  </WidgetColumnBar>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, unref, watch} from 'vue'
import {storeToRefs} from 'pinia'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {useLocaleStore} from '/@/store/modules/locale'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import WidgetColumnBar from './WidgetColumnBar.vue'
import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import ConnectionList from './ConnectionList.vue'
import SqlObjectList from './SqlObjectList.vue'
import PinnedObjectsList from './PinnedObjectsList'
import {findEngineDriver} from '/@/second/keeper-tools'
import {useConnectionInfo} from "/@/api/bridge";

export default defineComponent({
  name: "DatabaseWidget",
  props: {
    hidden: {
      type: Boolean as PropType<boolean>,
      default: false,
    }
  },
  components: {
    WidgetColumnBar,
    WidgetColumnBarItem,
    WidgetsInnerContainer,
    ConnectionList,
    SqlObjectList,
    PinnedObjectsList,
    ErrorInfo,
  },
  setup() {
    const bootstrap = useBootstrapStore()
    const {currentDatabase, extensions} = storeToRefs(bootstrap)
    const localeStore = useLocaleStore()
    const {pinnedDatabases, pinnedTables} = storeToRefs(localeStore)
    const database = computed(() => unref(currentDatabase)?.name)
    const conid = computed(() => unref(currentDatabase)?.connection._id)
    const singleDatabase = computed(() => unref(currentDatabase)?.connection?.singleDatabase)

    let connection = ref()
    watch(() => [conid.value, database.value], () => {
      useConnectionInfo({conid: unref(conid)}, connection)
    }, {
      immediate: true
    })

    const driver = computed(() => findEngineDriver(connection.value, extensions.value!))

    return {
      pinnedDatabases,
      pinnedTables,
      currentDatabase,
      conid,
      database,
      singleDatabase,
      driver
    }
  }
})
</script>
