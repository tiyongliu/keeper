import {computed,
  defineComponent,
  PropType,
  unref,
  toRefs
} from 'vue'
import {compact} from 'lodash-es'
import AppObjectListItem from '/@/second/appobj/AppObjectListItem.vue'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
import {Component} from "@vue/runtime-core";

import {IIsExpandable, IPinnedDatabasesItem} from '/@/second/types/standard.d'
export default defineComponent({
  name: "DatabaseWidget",
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
      groupFunc,
      filter,
      list,
      isExpandable,
      expandOnClick,
      passProps,
      subItemsComponent,
      expandIconFunc,
      module
    } = toRefs(props)


    const filtered = computed(() => {
      return !unref(groupFunc) ? (unref(list)!).filter(data => {
        const matcher = createMatcher && createMatcher(data);
        if (matcher && !matcher(filter)) return false;
        return true;
      }) : null
    })

    const childrenMatched = computed(() => {
      return !unref(groupFunc) ? (unref(list)!).filter(data => {
        const matcher = createChildMatcher && createChildMatcher(data)
        if (matcher && !matcher(filter)) return false;
        return true
      }) : null
    })

    const listGrouped = computed(() => {
      unref(groupFunc) ? compact(
        ((unref(list)!) || []).map(data => {
          const matcher = createMatcher && createMatcher(data);
          const isMatched = matcher && !matcher(filter) ? false : true;
        })
      ) : null
    })

    return () => (list.value || []).map(data => {
      return <AppObjectListItem
        isHidden={!(filtered.value as IPinnedDatabasesItem[]).includes(data)}
        module={unref(module)}
        subItemsComponent={unref(subItemsComponent)}
        data={unref(data)}
        isExpandable={unref(isExpandable)}

        expandIconFunc={unref(expandIconFunc)}
        isExpandedBySearch={(childrenMatched.value as IPinnedDatabasesItem[]).includes(data)}
        expandOnClick={unref(expandOnClick)}
        passProps={unref(passProps)}
      />
    })
  }
})
