import {computed,
  defineComponent,
  PropType,
  unref,
  toRefs
} from 'vue'
import {compact} from 'lodash-es'
import AppObjectListItem from '/@/second/appobj/AppObjectListItem.vue'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
import {IIsExpandable, IPinnedDatabasesItem} from '/@/second/types/standard.d'
import {Component} from "@vue/runtime-core";
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
      type: Function as PropType<(isExpanded: boolean) => string>,
    },
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    module: {
      type: [String, Object] as PropType<string | Component>,
    },
    subItemsComponent: {
      type: [String, Object] as PropType<string | Component>,
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
      console.log(list.value, `;list;`)
      console.log(subItemsComponent.value, `;subItemsComponent;`)
      return <AppObjectListItem
        isHidden={!(filtered.value as IPinnedDatabasesItem[]).includes(data)}
        module={unref(module)}
        subItemsComponent={unref(subItemsComponent)}
        expandOnClick={unref(expandOnClick)}
        data={unref(data)}
        isExpandable={unref(isExpandable)}

        expandIconFunc={unref(expandIconFunc)}
        filter={unref(filter)}
        isExpandedBySearch={(childrenMatched.value as IPinnedDatabasesItem[]).includes(data)}
        passProps={unref(passProps)}
      />
    })
  }
})
