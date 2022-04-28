<template>
  <template v-if="!isHidden">
    <component
      :is="module"
      :data="data"
      @click="handleExpand"
      @expand="handleExpandButton"
      :expandIcon="getExpandIcon(!isExpandedBySearch && expandable, subItemsComponent, isExpanded, expandIconFunc)"
      :passProps="passProps"
    />

    <span>isExpanded: {{ isExpanded }}</span> /
    <span>isExpandedBySearch: {{ isExpandedBySearch }}</span> /
    {{ subItemsComponent }}
    <div class="subitems" v-if="(isExpanded || isExpandedBySearch) && subItemsComponent">
      <component :is="subItemsComponent" :data="data" :filter="filter" :passProps="passProps" />
    </div>
  </template>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref, watch, ref, nextTick, onMounted, toRefs} from 'vue'
import {Component} from '@vue/runtime-core/dist/runtime-core'
import SubDatabaseList from './SubDatabaseList'
import ConnectionAppObject from './ConnectionAppObject'
import { plusExpandIcon } from '/@/second/icons/expandIcons';
import {IIsExpandable, IPinnedDatabasesItem} from '/@/second/types/standard.d'
import {getExpandIcon} from './module'
export default defineComponent({
  name: "AppObjectListItem",
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
    isHidden: {
      type: Boolean as unknown as PropType<boolean>,
    },
    isExpandedBySearch: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isExpandable: {
      type: Function as PropType<(IPinnedDatabasesItem) => any>,
      default: undefined
    },
    expandIconFunc: {
      type: Function as PropType<(isExpanded: boolean) => string>,
      default: plusExpandIcon
    },
    expandOnClick: {
      type: Boolean as PropType<boolean>,
    },
    module: {
      type: [Object, String] as PropType<string | Component>,
    },
    subItemsComponent: {
      type: [Object, String] as PropType<string | Component>,
    },
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    filter: {
      type: String as PropType<string>,
    }
  },
  components: {
    SubDatabaseList,
    ConnectionAppObject
  },
  setup(props) {
    const {data, isExpandable, subItemsComponent, expandOnClick, isExpandedBySearch} = toRefs(props)
    console.log(isExpandedBySearch, `isExpandedBySearch-isExpandedBySearch`)

    const isExpanded = ref(false)

    const expandable = computed(() => unref(data) && unref(isExpandable) && unref(isExpandable)!(data))

    async function handleExpand() {
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        isExpanded.value = !isExpanded.value
        console.log(isExpanded.value, `1111`)
      }
    }

    function handleExpandButton() {
      isExpanded.value = !isExpanded.value
    }

    watch(
      () => [unref(expandable), unref(isExpanded)],
      (watchExpandable, watchIsExpanded) => {
      if (!watchExpandable && watchIsExpanded) {
        isExpanded.value = false
      }
    })


    return {
      ...toRefs(props),
      expandable,
      isExpanded,
      handleExpand,
      handleExpandButton,
      getExpandIcon
    }
  }
})
</script>

<style lang="less">
.subitems {
  margin-left: 28px;
}
</style>
