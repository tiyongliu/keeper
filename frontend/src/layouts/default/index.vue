<template>
  <div class="not-supported">
    <div class="m-5 big-icon">
      <WarningOutlined />
    </div>
    <div class="m-3">Sorry, DbGate is not supported on mobile devices.</div>
    <div class="m-3">Please visit <a href="https://dbgate.org">DbGate web</a> for more info.</div>
  </div>

  <Layout :class="prefixCls" v-bind="lockEvents">
    <div class="iconbar svelte-1veekw4">
      <WidgetIconPanel @con="con"/>
    </div>
    <div class="statusbar svelte-1veekw4">iconbar</div>
<!--    <div class="statusbar svelte-1veekw4"><StatusBar></StatusBar></div>-->
    <div class="leftpanel svelte-1veekw4">
      <!--      <AppDarkModeToggle class="mx-auto" />-->

      <WidgetContainer :isShow="isShow"/>
    </div>
    <div class="tabs svelte-1veekw4">
<!--      <TabsPanel/>-->
      <LayoutHeader />
    </div>
    <div class="content svelte-1veekw4">content
    </div>
    <div class="horizontal-split-handle splitter svelte-1veekw4"
         v-splitterDrag="'clientX'"
         :resizeSplitter="(e) => cssVariable.setLeftPanelWidth(e.detail)">
    </div>
    <CurrentDropDownMenu />
    <div class="snackbar-container">snackbar-container</div>
  </Layout>
</template>

<script lang="ts">
  import { defineComponent, computed, unref, onMounted ,ref} from 'vue';
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

  //todo
  import { cssVariableStore } from "/@/store/modules/cssVariable"
  import WidgetContainer from '/@/second/widgets/WidgetContainer.vue'
  import TabsPanel from '/@/second/widgets/TabsPanel.vue'
  import StatusBar from '/@/second/widgets/StatusBar.vue'
  import {WarningOutlined} from '@ant-design/icons-vue'

  import WidgetIconPanel from '/@/second/widgets/WidgetIconPanel.vue'
  import CurrentDropDownMenu from '/@/second/modals/CurrentDropDownMenu'

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

      WidgetContainer,
      WarningOutlined,
      WidgetIconPanel,
      CurrentDropDownMenu,
      StatusBar,
      TabsPanel
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

      const cssVariable = cssVariableStore();
      onMounted(() => {
        cssVariable.subscribeCssVariable(cssVariable.getSelectedWidget,x => (x ? 1 : 0), '--dim-visible-left-panel')
        cssVariable.subscribeCssVariable(cssVariable.getLeftPanelWidth,x => `${x}px`, '--dim-left-panel-width')
        cssVariable.subscribeCssVariable(cssVariable.getVisibleTitleBar,x => (x ? 1 : 0), '--dim-visible-titlebar')

        //https://www.npmjs.com/package/resize-observer-polyfill
        // import ResizeObserver from 'resize-observer-polyfill';



      })
      const isShow = ref('database')
      

      return {
        getShowFullHeaderRef,
        getShowSidebar,
        prefixCls,
        getIsMobile,
        getIsMixSidebar,
        layoutClass,
        lockEvents,
        cssVariable,
        isShow
      };
    },
    methods:{
      con(val){
        this.isShow = val
      }
    }
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
body {
  font-family: -apple-system, BlinkMacSystemFont, Segoe WPC, Segoe UI, HelveticaNeue-Light, Ubuntu, Droid Sans,
  sans-serif;
  font-size: 14px;
  /* font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
    */
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

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

.icon-invisible {
  visibility: hidden;
}
.space-between {
  display: flex;
  justify-content: space-between;
}
.flex {
  display: flex;
}
.flexcol {
  display: flex;
  flex-direction: column;
}
.nowrap {
  white-space: nowrap;
}
.noselect {
  user-select: none;
}
.bold {
  font-weight: bold;
}
.flex1 {
  flex: 1;
}
.relative {
  position: relative;
}

.col-10 {
  flex-basis: 83.3333%;
  max-width: 83.3333%;
}
.col-9 {
  flex-basis: 75%;
  max-width: 75%;
}
.col-8 {
  flex-basis: 66.6667%;
  max-width: 66.6667%;
}
.col-7 {
  flex-basis: 58.3333%;
  max-width: 58.3333%;
}
.col-6 {
  flex-basis: 50%;
  max-width: 50%;
}
.col-5 {
  flex-basis: 41.6667%;
  max-width: 41.6667%;
}
.col-4 {
  flex-basis: 33.3333%;
  max-width: 33.3333%;
}
.col-3 {
  flex-basis: 25%;
  max-width: 25%;
}
.col-2 {
  flex-basis: 16.6666%;
  max-width: 16.6666%;
}

.largeFormMarker input[type='text'] {
  width: 100%;
  padding: 10px 10px;
  font-size: 14px;
  box-sizing: border-box;
  border-radius: 4px;
  border: 1px solid var(--theme-border);
}

.largeFormMarker input[type='password'] {
  width: 100%;
  padding: 10px 10px;
  font-size: 14px;
  box-sizing: border-box;
  border-radius: 4px;
}

.largeFormMarker select {
  width: 100%;
  padding: 10px 10px;
  font-size: 14px;
  box-sizing: border-box;
  border-radius: 4px;
}

body *::-webkit-scrollbar {
  height: 0.8em;
  width: 0.8em;
}
body *::-webkit-scrollbar-track {
  border-radius: 1px;
  background-color: var(--theme-bg-1);
}
body *::-webkit-scrollbar-corner {
  border-radius: 1px;
  background-color: var(--theme-bg-2);
}

body *::-webkit-scrollbar-thumb {
  border-radius: 1px;
  background-color: var(--theme-bg-3);
}

body *::-webkit-scrollbar-thumb:hover {
  background-color: var(--theme-bg-4);
}

input {
  background-color: var(--theme-bg-0);
  color: var(--theme-font-1);
  border: 1px solid var(--theme-border);
}

input[disabled] {
  background-color: var(--theme-bg-1);
}

select {
  background-color: var(--theme-bg-0);
  color: var(--theme-font-1);
  border: 1px solid var(--theme-border);
}

select[disabled] {
  background-color: var(--theme-bg-1);
}

textarea {
  background-color: var(--theme-bg-0);
  color: var(--theme-font-1);
  border: 1px solid var(--theme-border);
}
</style>
