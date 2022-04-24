<template>
    <WidgetColumnBar :hidden="hidden">
      <!--这个是上面数据库 及 db列表-->
      <WidgetColumnBarItem title="Connections" name="connections" height="35%" storageName="connectionsWidget">
        <ConnectionList />
      </WidgetColumnBarItem>
      <WidgetColumnBarItem
        title="Pinned"
        name="pinned"
        height="15%"
        storageName="pinnedItemsWidget"
        :skip="!pinnedDatabases?.length &&
        !pinnedTables.some(x => x.conid == currentDatabase.connection._id && x.database == currentDatabase?.name)"
      >
        <PinnedObjectsList />
      </WidgetColumnBarItem>

      <!--数据库 table 列表-->
      <template v-if="conid && (database || singleDatabase)">
        <WidgetColumnBarItem
          title="Tables, views, functions"
          name="dbObjects"
          storageName="dbObjectsWidget">
          <SqlObjectList :conid="conid" :database="database"/>
        </WidgetColumnBarItem>
      </template>
      <WidgetColumnBarItem v-else title="Database content" name="dbObjects" storageName="dbObjectsWidget">
        <WidgetsInnerContainer>
          <ErrorInfo message="Database not selected" icon="img alert"/>
        </WidgetsInnerContainer>
      </WidgetColumnBarItem>
    </WidgetColumnBar>
</template>

<script lang="ts">
  import {computed, defineComponent, PropType} from 'vue';
  import {dataBaseStore} from "/@/store/modules/dataBase"
  import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
  import WidgetColumnBar from './WidgetColumnBar.vue'
  import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  import ConnectionList from './ConnectionList.vue'
  import SqlObjectList from './SqlObjectList.vue'
  import PinnedObjectsList from './PinnedObjectsList'

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
      ErrorInfo
    },
    setup() {
      const dataBase = dataBaseStore()


      const pinnedDatabases = computed(() => dataBase.$state.pinnedDatabases)
      const pinnedTables = computed(() => dataBase.$state.pinnedTables)
      const conid = computed(() => dataBase.$state.currentDatabase?.connection._id)
      const currentDatabase = computed(() => dataBase.$state.currentDatabase)
      const database = computed(() => dataBase.$state.currentDatabase?.name)
      const singleDatabase = computed(() => dataBase.$state.currentDatabase?.connection?.singleDatabase)

      return {
        pinnedDatabases,
        pinnedTables,
        currentDatabase,
        conid,
        database,
        singleDatabase
      }
    }
  })
</script>

