<template>
<<<<<<< HEAD
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
=======
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
>>>>>>> 60c835d (整理)
export default defineComponent({
  name: "AppObjectListItem",
  props: {
    data: {
<<<<<<< HEAD
      type: Object as PropType<IPinnedDatabasesItem>,
=======
      type: Array as unknown as PropType<[]>,
      default: [],
>>>>>>> 60c835d (整理)
    },
    isHidden: {
      type: Boolean as unknown as PropType<boolean>,
    },
    isExpandedBySearch: {
<<<<<<< HEAD
      type: Boolean as PropType<boolean>,
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
=======
      type: Boolean as unknown as PropType<boolean>,
      default: false
    },
>>>>>>> 60c835d (整理)
    expandOnClick: {
      type: Boolean as PropType<boolean>,
    },
    module: {
<<<<<<< HEAD
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

    const isExpanded = ref(false)

    // const expandable = computed(() => unref(data) && unref(isExpandable) && isExpandable(data))
    const expandable = computed(() => unref(data) && unref(isExpandable))

    async function handleExpand() {
      alert(`handleExpand-handleExpand`)
      if (unref(subItemsComponent) && unref(expandOnClick)) {
        isExpanded.value = !isExpanded.value
        console.log(isExpanded.value, `1111`)
      }
    }

    function handleExpandButton() {
      isExpanded.value = !isExpanded.value
    }

    watch(() => [unref(expandable), unref(isExpanded)], (watchExpandable, watchIsExpanded) => {
      alert(`all`)
      console.log(`isExpanded`, unref(isExpanded))
      console.log(`isExpandedBySearch`, unref(isExpandedBySearch))
      console.log(`subItemsComponent`, unref(subItemsComponent))
      if (!watchExpandable && watchIsExpanded) {
        debugger
        isExpanded.value = false
      }
    })

    watch(() => [unref(isExpandedBySearch), unref(isExpandable)], (a, b) => {
      console.log(a, b)
    })

    return {
      ...toRefs(props),
      expandable,
      isExpanded,
      handleExpand,
      handleExpandButton,
      getExpandIcon
=======
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
>>>>>>> 60c835d (整理)
    }
  }
})
</script>

<style lang="less">
.subitems {
  margin-left: 28px;
}
</style>
