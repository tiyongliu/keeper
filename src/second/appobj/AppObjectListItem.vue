<template>
  <component
    v-if="!isHidden"
    :data="data"
    @click="handleExpand"
    :is="module"
  />
  <div class="subitems" v-if="subItemsComponent">
    <component :is="subItemsComponent" />
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref} from 'vue'
import {Component} from '@vue/runtime-core/dist/runtime-core'
export default defineComponent({
  name: "AppObjectListItem",
  props: {
    data: {
      type: Array as unknown as PropType<[]>,
      default: [],
    },
    isHidden: {
      type: Boolean as unknown as PropType<boolean>,
    },
    isExpandedBySearch: {
      type: Boolean as unknown as PropType<boolean>,
      default: false
    },
    expandOnClick: {
      type: Boolean as PropType<boolean>,
    },
    module: {
      type: String as PropType<Component>,
    },
    subItemsComponent: {
      type: String as PropType<Component>,
    }
  },
  setup(props) {
    // const dynamicList = computed(() => unref(props.list))
    // const currentComp = computed(() => unref(currentComp))
    async function handleExpand() {
      // if (subItemsComponent && expandOnClick) {
      //   await tick();
      //   isExpanded = !isExpanded;
      // }
    }

    function handleExpandButton() {

    }

    return {
      ...props,
      handleExpand,
      handleExpandButton,
    }
  }
})
</script>

<style lang="less">
.subitems {
  margin-left: 28px;
}
</style>
