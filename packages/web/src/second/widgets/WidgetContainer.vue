<template>
  <!--
  hidden={$selectedWidget != 'database'}
  -->
  <DatabaseWidget :hidden="isShow!=='database'"/>
  <FilesWidget :hidden="isShow!=='file'"/>
  <HistoryWidget :hidden="isShow!=='history'"/>
  <ArchiveWidget :hidden="isShow!=='archive'"/>
  <PluginsWidget :hidden="isShow!=='plugins'"/>
  <CellDataWidget :hidden="isShow!=='cell-data'"/>
  <AppWidget :hidden="isShow!=='app'"/>

  
</template>

<script lang="ts">
  import DatabaseWidget from './DatabaseWidget.vue'
  import CellDataWidget from './CellDataWidget.vue'
  import FilesWidget from './FilesWidget.vue'
  import HistoryWidget from './HistoryWidget.vue'
  import ArchiveWidget from './ArchiveWidget.vue'
  import PluginsWidget from './PluginsWidget.vue'
  import AppWidget from './AppWidget.vue'
  import { defineComponent, computed, unref, ref, onMounted,toRefs } from 'vue';
  export default defineComponent({
    props:{
      isShow:String
    },
    setup(props){
      const {isShow} = toRefs(props)
      return {
        isShow
      }
    },
    components: {
      DatabaseWidget,
      CellDataWidget,
      FilesWidget,
      HistoryWidget,
      ArchiveWidget,
      PluginsWidget,
      AppWidget
    }
  })
</script>
