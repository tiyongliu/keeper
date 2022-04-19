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
        :skip="[].length && [].some(x => x.conid == conid && x.database == '')"
      >
        <PinnedObjectsList />
      </WidgetColumnBarItem>

      <!--数据库 table 列表-->
      <WidgetColumnBarItem
        title="Connections"
        name="dbObjects"
        storageName="dbObjectsWidget"
      >
        <SqlObjectList />
      </WidgetColumnBarItem>
    </WidgetColumnBar>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, PropType} from 'vue';
  import { dataBaseStore } from "/@/store/modules/dataBase"
  import WidgetColumnBar from './WidgetColumnBar.vue'
  import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
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
      ConnectionList,
      SqlObjectList,
      PinnedObjectsList
    },
    setup() {
      const dataBase = dataBaseStore()
      return {
        pinnedDatabases: dataBase.$state.pinnedDatabases,
        pinnedTables: dataBase.$state.pinnedTables,
      }
    }
  })
</script>

<style scoped>

</style>
