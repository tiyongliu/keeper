import {Component, computed, defineComponent, PropType, toRaw, toRefs, unref} from 'vue'
import {compact, groupBy, keys} from 'lodash-es'
import AppObjectListItem from '/@/second/appobj/AppObjectListItem.vue'
import AppObjectGroup from '/@/second/appobj/AppObjectGroup.vue'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
import {IIsExpandable} from '/@/second/typings/types/standard.d'

export default defineComponent({
  name: "DatabaseWidget",
  props: {
    list: {
      type: Array as PropType<unknown[]>,
    },
    groupFunc: {
      type: Function as PropType<Function>,
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
      type: Object as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    module: {
      type: [String, Object] as PropType<string | Component>,
    },
    subItemsComponent: {
      type: [String, Object] as PropType<string | Component>,
    },
    disableContextMenu: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  setup(props) {
    const {
      groupFunc,
      list,
      filter,
      expandOnClick,
      passProps,
      expandIconFunc,
      isExpandable,
      disableContextMenu,
    } = toRefs(props)

    const module = toRaw(props.module)
    const subItemsComponent = toRaw(props.subItemsComponent)

    const filtered = computed(() => {
      return !unref(groupFunc) ? (unref(list)!).filter(data => {
        const matcher = createMatcher && createMatcher(data);
        if (matcher && !matcher(filter)) return false
        return true
      }) : null
    })

    const childrenMatched = computed(() => {
      return !unref(groupFunc) ? (unref(list)!).filter(data => {
        const matcher = createChildMatcher && createChildMatcher(data)
        if (matcher && !matcher(filter.value)) return false
        return true
      }) : null
    })


    function listGrouped() {
      return groupFunc.value ? compact(
        ((unref(list)!) || []).map(data => {
          const matcher = createMatcher && createMatcher(data);
          const isMatched = matcher && !matcher(filter.value) ? false : true;
          const group = groupFunc.value!(data)
          return {group, data, isMatched};
        })
      ) : null
    }

    const groups = computed<any>(() => unref(groupFunc) ? groupBy(listGrouped(), 'group') : null)

    function _AppObjectGroup() {
      return () => keys(unref(groups)).map(group => <AppObjectGroup
        group={group}
        module={unref(module)}
        items={unref(groups)![group]}
        expandIconFunc={unref(expandIconFunc)}
        isExpandable={unref(isExpandable)}
        subItemsComponent={unref(subItemsComponent)}
        groupFunc={unref(groupFunc)}
        disableContextMenu={unref(disableContextMenu)}
        filter={unref(filter)}
        passProps={unref(passProps)}
      />)
    }

    function _AppObjectListItem() {
      return () => (list.value || []).map(data => <AppObjectListItem
        isHidden={!(filtered.value)!.includes(data)}
        module={module}
        subItemsComponent={subItemsComponent}
        expandOnClick={unref(expandOnClick)}
        data={data as Record<string, any>}
        isExpandable={unref(isExpandable)}
        expandIconFunc={unref(expandIconFunc)}
        disableContextMenu={unref(disableContextMenu)}
        filter={unref(filter)}
        isExpandedBySearch={(childrenMatched.value)!.includes(data)}
        passProps={unref(passProps)}
      />)
    }

    return groupFunc.value ? _AppObjectGroup() : _AppObjectListItem()
  }
})
