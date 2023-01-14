<template>
  <div class="not-supported">
    <div class="m-5 big-icon">
      <WarningOutlined/>
    </div>
    <div class="m-3">Sorry, DbGate is not supported on mobile devices.</div>
    <div class="m-3">Please visit <a href="https://github.com/tiyongliu/keeper">keeper web</a> for more info.</div>
  </div>
  <div>
    <div class="iconbar">
      <WidgetIconPanel/>
    </div>
    <div class="statusbar">
      <StatusBar/>
    </div>
    <div v-if="selectedWidget" class="leftpanel">
      <WidgetContainer/>
    </div>
    <div class="tabs">
      <TabsPanel/>
    </div>
    <div class="content">
      <TabRegister/>
    </div>
    <div v-if="selectedWidget" class="horizontal-split-handle splitter"
         v-splitterDrag="'clientX'"
         :resizeSplitter="(e) => localeStore.setLeftPanelWidth(e.detail)">
    </div>
    <CurrentDropDownMenu/>
    <div class="snackbar-container">snackbar-container</div>
  </div>
</template>

<script lang="ts">
import {defineComponent, watch, ref, onMounted} from 'vue'
import {storeToRefs} from "pinia"
import {useLocaleStore} from '/@/store/modules/locale'
import {subscribeRecentDatabaseSwitch} from '/@/api/recentDatabaseSwitch'
import {subscribeCurrentDbByTab} from '/@/api/changeCurrentDbByTab'
import WidgetContainer from '/@/second/widgets/WidgetContainer.vue'
import TabsPanel from '/@/second/widgets/TabsPanel.vue'
import TabRegister from './TabRegister.vue'
import StatusBar from '/@/second/widgets/StatusBar.vue'
import {WarningOutlined} from '@ant-design/icons-vue'
import WidgetIconPanel from '/@/second/widgets/WidgetIconPanel.vue'
import CurrentDropDownMenu from '/@/second/modals/CurrentDropDownMenu'
import {debounce} from "lodash-es"
import bus from '/@/second/utility/bus'
export default defineComponent({
  name: 'DefaultLayout',
  components: {
    WidgetContainer,
    WarningOutlined,
    WidgetIconPanel,
    CurrentDropDownMenu,
    StatusBar,
    TabsPanel,
    TabRegister,
  },
  setup() {
    const excludeFirst = ref(false)
    const localeStore = useLocaleStore()
    const {selectedWidget, leftPanelWidth, visibleTitleBar} = storeToRefs(localeStore)
    window.addEventListener('resize', debounce(() => {
      if (excludeFirst.value) {
        bus.emitter.emit(bus.resize)
      }
    }, 300))

    onMounted(() => excludeFirst.value = true)

    watch(() => selectedWidget.value, () => {
      localeStore.setCssVariable(selectedWidget.value, x => (x ? 1 : 0), '--dim-visible-left-panel')
    }, {immediate: true})

    watch(() => leftPanelWidth.value, () => {
      localeStore.setCssVariable(leftPanelWidth.value, x => `${x}px`, '--dim-left-panel-width')
    }, {immediate: true})

    watch(() => visibleTitleBar.value, () => {
      localeStore.setCssVariable(visibleTitleBar.value, x => (x ? 1 : 0), '--dim-visible-titlebar')
    }, {immediate: true})

    subscribeCurrentDbByTab()
    subscribeRecentDatabaseSwitch()

    return {
      selectedWidget,
      localeStore
    };
  },
});
</script>

<style scoped>
.root {
  color: var(--theme-font-1);
}

.iconbar {
  position: fixed;
  display: flex;
  left: 0;
  top: var(--dim-header-top);
  bottom: var(--dim-statusbar-height);
  width: var(--dim-widget-icon-size);
  background: var(--theme-bg-inv-1);
}

.statusbar {
  position: fixed;
  background: var(--theme-bg-statusbar-inv);
  height: var(--dim-statusbar-height);
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
}

.leftpanel {
  position: fixed;
  top: var(--dim-header-top);
  left: var(--dim-widget-icon-size);
  bottom: var(--dim-statusbar-height);
  width: var(--dim-left-panel-width);
  background-color: var(--theme-bg-1);
  display: flex;
}

.tabs {
  position: fixed;
  top: var(--dim-header-top);
  left: var(--dim-content-left);
  height: var(--dim-tabs-panel-height);
  right: 0;
  background-color: var(--theme-bg-1);
  border-top: 1px solid var(--theme-border);
}

.content {
  position: fixed;
  top: var(--dim-content-top);
  left: var(--dim-content-left);
  bottom: var(--dim-statusbar-height);
  right: 0;
  background-color: var(--theme-bg-1);
}

.commads {
  position: fixed;
  top: var(--dim-header-top);
  left: var(--dim-widget-icon-size);
}

.toolbar {
  position: fixed;
  top: var(--dim-toolbar-top);
  height: var(--dim-toolbar-height);
  left: 0;
  right: 0;
  background: var(--theme-bg-1);
}

.splitter {
  position: absolute;
  top: var(--dim-header-top);
  bottom: var(--dim-statusbar-height);
  left: calc(var(--dim-widget-icon-size) + var(--dim-left-panel-width));
}

.snackbar-container {
  position: fixed;
  right: 0;
  bottom: var(--dim-statusbar-height);
}

.titlebar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--dim-titlebar-height);
}

.not-supported {
  display: none;
}

@media only screen and (max-width: 600px) {
  .dbgate-screen:not(.isElectron) {
    display: none;
  }

  .not-supported:not(.isElectron) {
    display: block;
  }
}

.not-supported {
  text-align: center;
}

.big-icon {
  font-size: 20pt;
}
</style>
