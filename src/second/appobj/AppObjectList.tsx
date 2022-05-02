import {
  computed,
  defineComponent,
  PropType,
  unref,
  toRefs,
  Component,
} from 'vue'
import {compact, keys, groupBy} from 'lodash-es'
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
  setup(props, {attrs}) {
    const {
      groupFunc,
      list,
      filter,
      expandOnClick,
      passProps,
      expandIconFunc,
      isExpandable,
      disableContextMenu,
      module,
      subItemsComponent,
    } = toRefs(props)

    const filtered = computed(() => {
      return !unref(groupFunc) ? (unref(list)!).filter(data => {
        const matcher =  createMatcher && createMatcher(data);
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

    // const listGrouped = computed(() => {
    //   groupFunc.value ? compact(
    //     ((unref(list)!) || []).map(data => {
    //       const matcher = createMatcher && createMatcher(data);
    //       const isMatched = matcher && !matcher(filter.value) ? false : true;
    //       const group = groupFunc.value!(data)
    //       return { group, data, isMatched };
    //     })
    //   ) : null
    // })


    function listGrouped() {
      return groupFunc.value ? compact(
        ((unref(list)!) || []).map(data => {
          const matcher = createMatcher && createMatcher(data);
          const isMatched = matcher && !matcher(filter.value) ? false : true;
          const group = groupFunc.value!(data)
          return { group, data, isMatched };
        })
      ) : null
    }

    function _AppObjectGroup() {
      const groups = groupBy(listGrouped(), 'group')

      return () => keys(groups).map(group => <AppObjectGroup
        group={group}
        module={unref(module)}
        items={groups[group]}
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
      console.log(`22222222222`, module.value)
      console.log(`22222222222`, subItemsComponent.value)



      return () => (list.value || []).map(data => <AppObjectListItem
        isHidden={!(filtered.value)!.includes(unref(data))}
        module={unref(module)}
        subItemsComponent={unref(subItemsComponent)}
        expandOnClick={unref(expandOnClick)}
        data={unref(data)}
        isExpandable={unref(isExpandable)}
        expandIconFunc={unref(expandIconFunc)}
        disableContextMenu={unref(disableContextMenu)}
        filter={unref(filter)}
        isExpandedBySearch={(childrenMatched.value)!.includes(unref(data))}
        passProps={unref(passProps)}
      />)
    }

    return groupFunc.value ? _AppObjectGroup() : _AppObjectListItem()
  }
})
