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
    <div class="subitems" v-if="(isExpanded || isExpandedBySearch) && subItemsComponent">
      <component :is="subItemsComponent" :data="data" :filter="filter" :passProps="passProps" />
    </div>
  </template>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref, watch, ref, nextTick, onMounted} from 'vue'
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
      type: Boolean as unknown as PropType<boolean>,
      default: false
    },
    isExpandable: {
      type: Function as PropType<IIsExpandable>,
      default: undefined
    },
    expandIconFunc: {
      type: [Function, Boolean] as PropType<(isExpanded: boolean) => string>,
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
    const {data, isExpandable, subItemsComponent, expandOnClick} = props

    const isExpanded = ref(false)

    const expandable = computed(() => unref(data) && unref(isExpandable) && isExpandable(data))

    async function handleExpand() {
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        await nextTick(() => {
          isExpanded.value = !isExpanded.value
        })
      }
    }

    function handleExpandButton() {
      isExpanded.value = !isExpanded.value
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
