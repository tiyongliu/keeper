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

    <span>isExpandedBySearch: {{ isExpandedBySearch }}</span>

    <template v-if="(isExpanded || isExpandedBySearch) && subItemsComponent">
      <div class="subitems">
        <component :is="subItemsComponent" :data="data" :filter="filter" :passProps="passProps"/>
      </div>
    </template>
  </template>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'
import {Component} from '@vue/runtime-core/dist/runtime-core'
import {plusExpandIcon} from '/@/second/icons/expandIcons';
import {IPinnedDatabasesItem} from '/@/second/types/standard.d'
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
      type: Object as unknown as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
    },
    filter: {
      type: String as PropType<string>,
    }
  },
  setup(props) {
    const {data, isExpandable, subItemsComponent, expandOnClick} = toRefs(props)
    const isExpanded = ref(false)

    const expandable = computed(() => unref(data) && unref(isExpandable) && unref(isExpandable)!(data))

    async function handleExpand() {
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        isExpanded.value = !isExpanded.value
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
