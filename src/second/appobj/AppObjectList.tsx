import {computed, defineComponent, onMounted, PropType, unref, watch} from 'vue'
import {compact} from 'lodash-es'
import AppObjectListItem from './AppObjectListItem.vue'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
import {Component} from "@vue/runtime-core";

import {IIsExpandable, IPinnedDatabasesItem} from '/@/second/types/IStore.d'
export default defineComponent({
  name: "DatabaseWidget",
  props: {
    list: {
      type: Array as unknown as PropType<IPinnedDatabasesItem[]>,
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
    // const dynamicList = computed(() => unref(props.list))
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

    return () => (
      (unref(list)!).map(data => {
        return <AppObjectListItem
          isHidden={!(filtered.value as IPinnedDatabasesItem[]).includes(data)}
          module={unref(module)}
          subItemsComponent={unref(subItemsComponent)}
          data={unref(data)}
          isExpandable={isExpandable}

          expandIconFunc={unref(expandIconFunc)}
          isExpandedBySearch={(childrenMatched.value as IPinnedDatabasesItem[]).includes(data)}
          expandOnClick={unref(expandOnClick)}
          passProps={unref(passProps)}
        />
      })
    )
  }
})
