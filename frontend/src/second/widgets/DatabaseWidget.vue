<template>
  <WidgetColumnBar :hidden="hidden">
    <!--这个是上面数据库 及 db列表-->
    <WidgetColumnBarItem
      title="Connections"
      name="connections"
      height="35%"
      storageName="connectionsWidget">
      <ConnectionList/>
    </WidgetColumnBarItem>
    <WidgetColumnBarItem
      title="Pinned"
      name="pinned"
      height="15%"
      storageName="pinnedItemsWidget"
      :skip="!pinnedDatabases.length &&
      !pinnedTables.some(x =>
      currentDatabase && x.conid == currentDatabase.connection._id && x.database == currentDatabase?.name)">
      <PinnedObjectsList/>
    </WidgetColumnBarItem>
    <!--数据库 table 列表-->
    <WidgetColumnBarItem
      v-if="conid && (database || singleDatabase)"
      :title="driver && Array.isArray(driver?.databaseEngineTypes) && driver?.databaseEngineTypes?.includes('document')
      ? 'Collections'
      : 'Tables, views, functions'"
      name="dbObjects"
      storageName="dbObjectsWidget">
      <SqlObjectList :conid="conid" :database="database"/>
    </WidgetColumnBarItem>
    <WidgetColumnBarItem
      v-else
      title="Database content"
      name="dbObjects"
      storageName="dbObjectsWidget">
      <WidgetsInnerContainer>
        <ErrorInfo message="Database not selected" icon="img alert"/>
      </WidgetsInnerContainer>
    </WidgetColumnBarItem>
  </WidgetColumnBar>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, toRef, unref} from 'vue'
import {storeToRefs} from 'pinia'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {useClusterApiStore} from '/@/store/modules/clusterApi'
import {useLocaleStore} from '/@/store/modules/locale'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import WidgetColumnBar from './WidgetColumnBar.vue'
import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import ConnectionList from './ConnectionList.vue'
import SqlObjectList from './SqlObjectList.vue'
import PinnedObjectsList from './PinnedObjectsList'
import {findEngineDriver} from '/@/second/keeper-tools'

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
  setup(props) {
    const bootstrap = useBootstrapStore()
    const {currentDatabase, extensions} = storeToRefs(bootstrap)
    const localeStore = useLocaleStore()
    const {pinnedDatabases, pinnedTables} = storeToRefs(localeStore)

    const clusterApi = useClusterApiStore()
    const {connection} = storeToRefs(clusterApi)

    const database = computed(() => unref(currentDatabase)?.name)
    const conid = computed(() =>
      (unref(currentDatabase) && unref(currentDatabase)!.connection)
      ? unref(currentDatabase)?.connection._id : null)
    const driver = computed(() => extensions.value ? findEngineDriver(connection.value, extensions.value) : null)
    const singleDatabase = computed(() => unref(currentDatabase)?.connection?.singleDatabase)

    return {
      hidden: toRef(props, 'hidden'),
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
