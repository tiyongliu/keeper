<template>
  <AppObjectListItem
    v-for="data in list"
    :isHidden="!filtered.includes(data)"
    :module="module"
    :subItemsComponent="subItemsComponent"
    :data="data"
    :isExpandable="isExpandable"
    :expandIconFunc="expandIconFunc"
    :isExpandedBySearch="childrenMatched.includes(data)"
    :expandOnClick="expandOnClick"
    :passProps="passProps"
  />
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, unref, watch} from 'vue'
import {compact} from 'lodash-es'
import AppObjectListItem from './AppObjectListItem.vue'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
import {IIsExpandable, IPinnedDatabasesItem} from "/@/second/types/IStore";
import {Component} from "@vue/runtime-core";
import {dataBaseStore} from "/@/store/modules/dataBase";

export default defineComponent({
  name: "AppObjectList",
  components: {AppObjectListItem},
  props: {
    list: {
      type: Array as PropType<IPinnedDatabasesItem[]>,
    },
    groupFunc: {
      type: String as PropType<string>,
    },
    expandOnClick: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isExpandable: {
      type: Function as PropType<IIsExpandable>
    },
    filter: {
      type: String as PropType<string>,
    },
    expandIconFunc: {
      type: [Function, Boolean] as PropType<(isExpanded: boolean) => string>,
    },
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    module: {
      type: [Object, String] as PropType<string | Component>,
    },
    subItemsComponent: {
      type: [Object, String] as PropType<string | Component>,
    },
  },
  setup(props) {
    const {
      list,
      groupFunc,
      filter,
      isExpandable,
      expandOnClick,
      passProps,
      subItemsComponent,
      expandIconFunc,
      module
    } = props

    const filtered = computed(() => {
      return !unref(groupFunc) ? (list!).filter(data => {
        const matcher = createMatcher && createMatcher(data);
        if (matcher && !matcher(filter)) return false;
        return true;
      }) : null
    })

    const childrenMatched = computed(() => {
      return !unref(groupFunc) ? (list!).filter(data => {
        const matcher = createChildMatcher && createChildMatcher(data)
        if (matcher && !matcher(filter)) return false;
        return true
      }) : null
    })

    const listGrouped = computed(() => {
      unref(groupFunc) ? compact(
        ((list!) || []).map(data => {
          const matcher = createMatcher && createMatcher(data);
          const isMatched = matcher && !matcher(filter) ? false : true;
        })
      ) : null
    })

    const dataBase = dataBaseStore()

    watch(() => unref(dataBase.$state.pinnedDatabases), () => {
      // console.log(unref(list), ` unref(list)          unref(list)`)
    })

    watch(() => props.list, () => {
      console.log(props.list, ` unref(list)          unref(list)`)
    })


    onMounted(() => {
      // console.log(props.list, 'onMounted')
    })

    return {
      ...props,
      filtered,
      childrenMatched,
    }
  }
})
</script>

<style scoped>

</style>
