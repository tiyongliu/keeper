<template>
  <component
    v-if="!isHidden"
    :data="data"
    @click="handleExpand"
    :is="module"
  />
  <div class="subitems" v-if="subItemsComponent">
    <component :is="subItemsComponent" :data="data"/>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref, onMounted} from 'vue'
import {Component} from '@vue/runtime-core/dist/runtime-core'

import SubDatabaseList from './SubDatabaseList'
import ConnectionAppObject from './ConnectionAppObject'
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
      type: [Object, String] as PropType<string | Component>,
    },
    subItemsComponent: {
      type: [Object, String] as PropType<string | Component>,
    }
  },
  components: {
    SubDatabaseList,
    ConnectionAppObject
  },
  setup(props) {
    // const dynamicList = computed(() => unref(props.list))
    // const currentComp = computed(() => unref(currentComp))
    async function handleExpand() {
      alert(`handleExpand`)
      // if (subItemsComponent && expandOnClick) {
      //   await tick();
      //   isExpanded = !isExpanded;
      // }
    }

    function handleExpandButton() {

    }

    onMounted(() => {
      console.log(`onMounted`, props.data)
    })
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
