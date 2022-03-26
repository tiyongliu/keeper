<template>
  <div class="not-supported">
    <div class="m-5 big-icon">
      <warning-outlined />
    </div>
    <div class="m-3">Sorry, DbGate is not supported on mobile devices.</div>
    <div class="m-3">Please visit <a href="https://dbgate.org">DbGate web</a> for more info.</div>
  </div>

  <Layout :class="prefixCls" class="dbgate-screen" v-bind="lockEvents">
    <div class="iconbar svelte-1veekw4">iconbar</div>
    <div class="statusbar svelte-1veekw4">iconbar</div>
    <div class="leftpanel svelte-1veekw4">
<!--      <AppDarkModeToggle class="mx-auto" />-->

      <WidgetContainer />
    </div>
    <div class="tabs svelte-1veekw4">
      <LayoutHeader />
    </div>
    <div class="content svelte-1veekw4">content
    </div>
    <div class="horizontal-split-handle splitter svelte-1veekw4"
         v-splitterDrag="'clientX'"
         :resizeSplitter="(e) => varStore.setLeftPanelWidth(e.detail)">
    </div>
    <div class="snackbar-container svelte-1veekw4">snackbar-container</div>
  </Layout>
</template>

<script lang="ts">
import { defineComponent, computed, unref, ref, onMounted } from 'vue';
import { Layout } from 'ant-design-vue';
import { createAsyncComponent } from '/@/utils/factory/createAsyncComponent';

import LayoutHeader from './header/index.vue';
import LayoutContent from './content/index.vue';
import LayoutSideBar from './sider/index.vue';
import LayoutMultipleHeader from './header/MultipleHeader.vue';

import { useHeaderSetting } from '/@/hooks/setting/useHeaderSetting';
import { useMenuSetting } from '/@/hooks/setting/useMenuSetting';
import { useDesign } from '/@/hooks/web/useDesign';
import { useLockPage } from '/@/hooks/web/useLockPage';

import { useAppInject } from '/@/hooks/web/useAppInject';

import { AppDarkModeToggle } from '/@/components/Application';


//todo
import { variableStore } from "/@/store/modules/variable";
import WidgetContainer from '/@/second/widgets/WidgetContainer.vue'
import {WarningOutlined} from '@ant-design/icons-vue'
export default defineComponent({
  name: 'DefaultLayout',
  components: {
    LayoutFeatures: createAsyncComponent(() => import('/@/layouts/default/feature/index.vue')),
    LayoutFooter: createAsyncComponent(() => import('/@/layouts/default/footer/index.vue')),
    LayoutHeader,
    LayoutContent,
    LayoutSideBar,
    LayoutMultipleHeader,
    Layout,
    AppDarkModeToggle,


    WidgetContainer,
    WarningOutlined
  },
  setup() {
    const { prefixCls } = useDesign('default-layout');
    const { getIsMobile } = useAppInject();
    const { getShowFullHeaderRef } = useHeaderSetting();
    const { getShowSidebar, getIsMixSidebar, getShowMenu } = useMenuSetting();

    // Create a lock screen monitor
    const lockEvents = useLockPage();

    const layoutClass = computed(() => {
      let cls: string[] = ['ant-layout'];
      if (unref(getIsMixSidebar) || unref(getShowMenu)) {
        cls.push('ant-layout-has-sider');
      }
      return cls;
    });

    const varStore = variableStore();

    onMounted(() => {
      varStore.subscribeCssVariable(varStore.getSelectedWidget,x => (x ? 1 : 0), '--dim-visible-left-panel')
      varStore.subscribeCssVariable(varStore.getLeftPanelWidth,x => `${x}px`, '--dim-left-panel-width')
      varStore.subscribeCssVariable(varStore.getVisibleTitleBar,x => (x ? 1 : 0), '--dim-visible-titlebar')

      //https://www.npmjs.com/package/resize-observer-polyfill
      // import ResizeObserver from 'resize-observer-polyfill';

    })

    return {
      getShowFullHeaderRef,
      getShowSidebar,
      prefixCls,
      getIsMobile,
      getIsMixSidebar,
      layoutClass,
      lockEvents,
      varStore
    };
  },
});
</script>
<style lang="less">
@prefix-cls: ~'@{namespace}-default-layout';

.@{prefix-cls} {
  display: flex;
  width: 100%;
  min-height: 100%;
  background-color: @content-bg;
  flex-direction: column;

  > .ant-layout {
    min-height: 100%;
  }

  &-main {
    width: 100%;
    margin-left: 1px;
  }
}
</style>


<style lang="less">
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

<style>
.horizontal-split-handle {
  background-color: var(--theme-border);
  width: var(--dim-splitter-thickness);
  cursor: col-resize;
}
.horizontal-split-handle:hover {
  background-color: var(--theme-bg-2);
}

.vertical-split-handle {
  background-color: var(--theme-border);
  height: var(--dim-splitter-thickness);
  cursor: row-resize;
}
.vertical-split-handle:hover {
  background-color: var(--theme-bg-2);
}

</style>
