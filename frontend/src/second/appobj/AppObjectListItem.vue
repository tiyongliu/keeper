<template>
  <template v-if="!isHidden">
    <component
      :is="module"
      :data="data"
      @click="handleExpand"
      @expand="handleExpandButton"
      :expandIcon="getExpandIcon(!isExpandedBySearch && expandable, subItemsComponent, isExpanded, expandIconFunc)"
      :disableContextMenu="disableContextMenu"
      :passProps="passProps"
      v-bind="$attrs"
    />

    <div class="subitems" v-if="(isExpanded || isExpandedBySearch) && subItemsComponent">
      <component
        :is="subItemsComponent"
        :data="data"
        :filter="filter"
        :passProps="passProps"
        v-bind="$attrs"/>
    </div>
  </template>
</template>

<script lang="ts">
import {Component, computed, defineComponent, PropType, ref, toRaw, toRefs, unref, watch} from 'vue'
import {plusExpandIcon} from '/@/second/icons/expandIcons';
import {getExpandIcon} from './module'

export default defineComponent({
  name: "AppObjectListItem",
  props: {
    data: {
      type: Object as PropType<Record<string, any>>,
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
      type: Object as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
    },
    filter: {
      type: String as PropType<string>,
    },
    disableContextMenu: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    getIsExpanded: {
      type: Function as PropType<(data: any) => boolean>,
    },
    setIsExpanded: {
      type: Function as PropType<(data: any, isExpanded: boolean) => void>,
    }
  },
  setup(props) {
    const {
      data,
      isHidden,
      isExpandedBySearch,
      isExpandable,
      expandIconFunc,
      expandOnClick,
      passProps,
      filter,
      getIsExpanded,
      setIsExpanded,
      disableContextMenu
    } = toRefs(props)
    const module = toRaw(props.module)
    const subItemsComponent = toRaw(props.subItemsComponent)
    const isExpandedCore = ref(false)

    const expandable = computed(() => {
      return data.value && isExpandable.value && isExpandable.value(data.value)
    })

    const isExpanded = computed(() => {
      return expandable.value ? (getIsExpanded.value && setIsExpanded.value ? getIsExpanded.value(data.value) : isExpandedCore.value) : false
    })

    async function handleExpand() {
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        handleExpandButton()
      }
    }

    function handleExpandButton() {
      if (getIsExpanded.value && setIsExpanded.value) {
        setIsExpanded.value(data.value, !isExpanded.value)
      } else {
        isExpandedCore.value = !isExpandedCore.value
      }
    }

    return {
      data,
      isHidden,
      isExpandedBySearch,
      isExpandable,
      expandIconFunc,
      expandOnClick,
      passProps,
      filter,
      disableContextMenu,
      module,
      subItemsComponent,
      expandable,
      isExpanded,
      handleExpand,
      handleExpandButton,
      getExpandIcon,
    }
  }
})
</script>

<style lang="less" scoped>
.subitems {
  margin-left: 28px;
}
</style>
