<template>
  <Layout :class="prefixCls" v-bind="lockEvents">
    <div class="leftpanel" :style="{width: `${width}px`}">
      <AppDarkModeToggle class="mx-auto" />
    </div>
    <div class="tabs">
      <LayoutHeader />
    </div>
    <div class="content">3</div>
    <div class="horizontal-split-handle splitter"
         @mousedown.prevent="writableWithStorage">4</div>
  </Layout>
</template>

<script lang="ts">
  import { defineComponent, computed, unref, ref } from 'vue';
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
    },
    setup() {
      const width = ref(230)

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


      const writableWithStorage = (event: MouseEvent) => {
        console.log(event.detail)
        console.log('获取当前坐标', event.clientX)
        document.onmousemove = function (e) {
          // console.log('获取移动的')
          // moveHandle(e.clientY)
          return false
        }

        document.onmouseup = function() {
          document.onmousemove = null
          document.onmouseup = null
        }
        return false
      }

      const moveHandle = (nowClient) => {
        console.log(nowClient)
        //改动窗口大小
        //https://blog.csdn.net/heixiuheixiu666/article/details/107609128
        // console.log(document.querySelectorAll('.leftpanel')[0].clientWidth)
        // document.querySelectorAll('.leftpanel')[0].clientWidth = nowClient
      }


      return {
        getShowFullHeaderRef,
        getShowSidebar,
        prefixCls,
        getIsMobile,
        getIsMixSidebar,
        layoutClass,
        lockEvents,

        width,
        writableWithStorage
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
  .leftpanel {
    position: fixed;
    top: var(--dim-header-top);
    left: var(--dim-widget-icon-size);
    bottom: var(--dim-statusbar-height);
    /*width: var(--dim-left-panel-width);*/
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

  .splitter {
    position: absolute;
    top: var(--dim-header-top);
    bottom: var(--dim-statusbar-height);
    left: calc(var(--dim-widget-icon-size) + var(--dim-left-panel-width));
  }
</style>

<style lang="less">
  .horizontal-split-handle {
    background-color: var(--theme-border);
    width: var(--dim-splitter-thickness);
    cursor: col-resize;
  }
  .horizontal-split-handle:hover {
    background-color: var(--theme-bg-2);
  }



  .horizontal-split-handle {
    background-color:  #ccc;
    width: var(--dim-splitter-thickness);
    cursor: col-resize;
  }
  .horizontal-split-handle:hover {
    background-color: #d4d4d4;
  }
</style>
