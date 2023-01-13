<template>
  <div class="main" :class="flex1 && 'flex1'">
    <div class="tabs">
      <div class="tab-item" v-for="(tab, index) in tabs">
        <div class="tab-item" :class="value == index && 'selected'" @click="targetIndex(index)">
          <span class="ml-2">
            {{ tab.label }}
          </span>
        </div>
      </div>

      <DropDownButton v-if="menu" :menu="menu"/>
    </div>

    <div class="content-container">
      <div
        class="container"
        v-for="(tab, index) in tabs"
        :class="[isInline && 'isInline', index == value && 'tabVisible']"
        :style="`max-width: ${containerMaxWidth}`">
        <component :is="tab.component" v-bing="{...tab.props}"
                   :tabControlHiddenTab="index != value"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {Component, defineComponent, PropType, toRefs} from "vue"
import DropDownButton from '/@/second/buttons/DropDownButton'
import {compact} from "lodash-es";
import {Tabs} from "ant-design-vue";

const TabPane = Tabs.TabPane

interface TabDef {
  label: string;
  slot?: number;
  component?: string | Component;
  props?: any;
}

export default defineComponent({
  name: 'TabControl',
  components: {
    DropDownButton,
    [Tabs.name]: Tabs,
    [TabPane.name]: TabPane,
  },
  props: {
    isInline: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    tabs: {
      type: Array as PropType<TabDef[]>
    },
    value: {
      type: Number as PropType<number>,
      default: 0
    },
    menu: {
      type: Array as unknown as PropType<[]>,
    },
    containerMaxWidth: {
      type: String as PropType<string>
    },
    flex1: {
      type: Boolean as PropType<boolean>,
      default: true
    }
  },
  setup(props) {
    const {isInline, tabs, value, menu, containerMaxWidth, flex1} = toRefs(props)

    function targetIndex(index: number) {
      value.value = index
    }

    return {
      isInline,
      menu,
      containerMaxWidth,
      flex1,
      tabs: compact(tabs.value),
      targetIndex,
    }
  }
})
</script>

<style scoped>
.main {
  display: flex;
  flex-direction: column;
}

.main.flex1 {
  flex: 1;
}

.tabs {
  display: flex;
  height: var(--dim-tabs-height);
  right: 0;
  background-color: var(--theme-bg-2);
}

.tab-item {
  border-right: 1px solid var(--theme-border);
  padding-left: 15px;
  padding-right: 15px;
  display: flex;
  align-items: center;
  cursor: pointer;
}

/* .tab-item:hover {
  color: ${props => props.theme.tabs_font_hover};
} */
.tab-item.selected {
  background-color: var(--theme-bg-1);
}

.content-container {
  flex: 1;
  position: relative;
}

.container:not(.isInline) {
  position: absolute;
  display: flex;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

.container:not(.tabVisible):not(.isInline) {
  visibility: hidden;
}

.container.isInline:not(.tabVisible) {
  display: none;
}
</style>


<!--<template>-->
<!--  <div>-->
<!--    <div>-->
<!--      <a-tabs default-active-key="1">-->
<!--        <a-tab-pane :key="`${index + 1}`" :tab="tab.label" v-for="(tab, index) in tabs">-->
<!--          <component :is="tab.component"/>-->
<!--        </a-tab-pane>-->
<!--      </a-tabs>-->
<!--    </div>-->

<!--    <div></div>-->
<!--  </div>-->
<!--</template>-->

<!--<script lang="ts">-->
<!--import {Component, defineComponent} from "vue";-->
<!--import {compact} from "lodash-es";-->
<!--import { Tabs} from "ant-design-vue";-->
<!--const TabPane = Tabs.TabPane-->
<!--interface TabDef {-->
<!--  label: string;-->
<!--  slot?: number;-->
<!--  component?: string | Component;-->
<!--  props?: any;-->
<!--}-->
<!--export default defineComponent({-->
<!--  name: 'TabControl',-->
<!--  components: {-->
<!--    [Tabs.name]: Tabs,-->
<!--    [TabPane.name]: TabPane,-->
<!--  },-->
<!--  props: {-->
<!--    isInline: {-->
<!--      type: Boolean as PropType<boolean>,-->
<!--      default: false-->
<!--    },-->
<!--    tabs: {-->
<!--      type: Array as PropType<TabDef[]>-->
<!--    }-->
<!--  },-->
<!--  setup(props) {-->
<!--    return {-->
<!--      tabs: compact(props.tabs)-->
<!--    }-->
<!--  }-->
<!--})-->
<!--</script>-->

<!--<style scoped>-->

<!--</style>-->
