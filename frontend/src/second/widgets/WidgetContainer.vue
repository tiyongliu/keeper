<template>
  <DatabaseWidget :hidden="selectedWidget != 'database'"/>
  <FilesWidget v-if="selectedWidget == 'file'"/>
  <HistoryWidget v-if="selectedWidget == 'history'"/>
  <ArchiveWidget  v-if="selectedWidget == 'archive'"/>
  <PluginsWidget v-if="selectedWidget == 'plugins'"/>
  <CellDataWidget v-if="selectedWidget == 'cell-data'"/>
  <AppWidget v-if="selectedWidget == 'app'"/>
</template>

<script lang="ts">
import DatabaseWidget from './DatabaseWidget.vue'
import CellDataWidget from './CellDataWidget.vue'
import FilesWidget from './FilesWidget.vue'
import HistoryWidget from './HistoryWidget.vue'
import ArchiveWidget from './ArchiveWidget.vue'
import PluginsWidget from './PluginsWidget.vue'
import AppWidget from './AppWidget.vue'
import {defineComponent} from 'vue';
import {useLocaleStore} from '/@/store/modules/locale'
import {storeToRefs} from "pinia";
export default defineComponent({
  setup() {
    const localeStore = useLocaleStore()
    const {selectedWidget} = storeToRefs(localeStore)
    return {
      selectedWidget
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
