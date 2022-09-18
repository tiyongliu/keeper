<template>
  <div class="main">
    <div class="wrapper" :class="{'selected':item.name === selectedWidget}"
         v-for="(item, index) in widgets" :key="index"
         @click="handleChangeWidget(item.name)">
      <FontIcon :icon="item.icon" :title="item.title"/>
    </div>

    <div class="flex1">&nbsp;</div>

    <div class="wrapper" @click="handleSettingsMenu" ref="domSettings">
      <FontIcon icon="icon settings"/>
    </div>
  </div>
</template>

<script lang="ts">
import {reactive, ref} from 'vue';
import {storeToRefs} from 'pinia'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useLocaleStore} from '/@/store/modules/locale'
import {useBootstrapStore} from '/@/store/modules/bootstrap'

export default {
  setup() {
    const widgets = reactive([
      {
        icon: 'icon database',
        name: 'database',
        title: 'Database connections',
      },
      // {
      //   icon: 'fa-table',
      //   name: 'table',
      // },
      {
        icon: 'icon file',
        name: 'file',
        title: 'Favorites & Saved files',
      },
      {
        icon: 'icon history',
        name: 'history',
        title: 'Query history & Closed tabs',
      },
      {
        icon: 'icon archive',
        name: 'archive',
        title: 'Archive (saved tabular data)',
      },
      {
        icon: 'icon plugin',
        name: 'plugins',
        title: 'Extensions & Plugins',
      },
      {
        icon: 'icon cell-data',
        name: 'cell-data',
        title: 'Selected cell data detail view',
      },
      {
        icon: 'icon app',
        name: 'app',
        title: 'Application layers',
      },
      // {
      //   icon: 'icon settings',
      //   name: 'settings',
      // },
      // {
      //   icon: 'fa-check',
      //   name: 'settings',
      // },
    ])

    const localeStore = useLocaleStore()
    const bootstrap = useBootstrapStore()
    const {selectedWidget} = storeToRefs(localeStore)
    const domSettings = ref<Nullable<HTMLElement>>(null)

    function handleSettingsMenu() {
      const rect = domSettings.value!.getBoundingClientRect();
      const left = rect.right
      const top = rect.bottom
      const items = [{command: 'settings.show'}, {command: 'theme.changeTheme'}, {command: 'settings.commands'}]
      bootstrap.subscribeCurrentDropDownMenu({left, top, items})
    }

    function handleChangeWidget(name) {
      localeStore.setSelectedWidget(name === selectedWidget.value ? null : name)
    }

    return {
      widgets,
      selectedWidget,
      domSettings,
      handleSettingsMenu,
      handleChangeWidget
    };
  },
  components: {FontIcon},
};
</script>

<style scoped>
.wrapper {
  font-size: 23pt;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--theme-font-inv-2);
}

.wrapper:hover {
  color: var(--theme-font-inv-1);
}

.wrapper.selected {
  color: var(--theme-font-inv-1);
  background: var(--theme-bg-inv-3);
}

.main {
  display: flex;
  flex: 1;
  flex-direction: column;
}
</style>


