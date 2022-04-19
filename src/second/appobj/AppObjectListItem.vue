<template>
  <template v-if="!isHidden">
    <component
      :is="module"
      :data="data"
      @click="handleExpand"
      :expandIcon="getExpandIcon(!isExpandedBySearch && expandable, subItemsComponent, isExpanded, expandIconFunc)"
      :passProps="passProps"
    />
    <div class="subitems" v-if="(isExpanded || isExpandedBySearch) && subItemsComponent">
      <component :is="subItemsComponent" :data="data" :filter="filter" :passProps="passProps" />
    </div>
  </template>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref, watch, ref, nextTick} from 'vue'
import {Component} from '@vue/runtime-core/dist/runtime-core'
import SubDatabaseList from './SubDatabaseList'
import ConnectionAppObject from './ConnectionAppObject'
import { plusExpandIcon } from '/@/second/icons/expandIcons';
import {IIsExpandable} from '/@/second/types/IStore.d'
import {getExpandIcon} from './module'
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
    isExpandable: {
      type: Function as PropType<IIsExpandable>,
      default: undefined
    },
    expandIconFunc: {
      type: Function as PropType<(isExpanded: boolean) => string>,
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
  },
  components: {
    SubDatabaseList,
    ConnectionAppObject
  },
  setup(props) {
    const {data, isExpandable, subItemsComponent, expandOnClick} = props

    const isExpanded = ref(false)

    const expandable = computed(() => unref(data) && unref(isExpandable) && isExpandable(data))

    // const dynamicList = computed(() => unref(props.list))
    // const currentComp = computed(() => unref(currentComp))
    async function handleExpand() {
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        await nextTick(() => {
          isExpanded.value = !isExpanded.value
        })
      }
    }

    function handleExpandButton() {

    }

    const handle = () => {
      if (unref(expandable) && unref(isExpandable)) {
        isExpanded.value = false
      }
    }

    watch(() => [unref(expandable), unref(isExpanded)], handle)

    return {
      ...props,
      expandable,
      isExpanded,
      handleExpand,
      handleExpandButton,
      getExpandIcon,
    }
  }
})
</script>

<style lang="less">
.subitems {
  margin-left: 28px;
}
</style>
