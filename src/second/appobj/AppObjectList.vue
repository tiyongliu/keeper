<template>
  <AppObjectListItem
    :isHidden="!filtered.includes(data)"
    :data="data"
    :module="ConnectionAppObject"
    :subItemsComponent="SubDatabaseList"
    :isExpandable="isExpandable"
    :isExpandedBySearch="childrenMatched.includes(data)"
  />
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, unref, watch} from 'vue'
import {compact} from 'lodash-es'
import AppObjectListItem from './AppObjectListItem.vue'
import SubDatabaseList from './SubDatabaseList'
import ConnectionAppObject from './ConnectionAppObject'
import {createChildMatcher, createMatcher} from './ConnectionAppObject'
type fn = (data: {_id: string, singleDatabase: boolean}) => boolean
  export default defineComponent({
    name: "DatabaseWidget",
    components: {AppObjectListItem},
    props: {
      list: {
        type: Array as unknown as PropType<[]>,
      },
      groupFunc: {
        type: String as PropType<string>,
      },
      expandOnClick: {
        type: Boolean as PropType<boolean>,
        default: false
      },
      isExpandable: {
        type: Function as PropType<fn>
      },
      filter: {
        type: String as PropType<string>,
      }
    },
    setup(props) {
      const {list, groupFunc, filter, isExpandable} = props

      const filtered = computed(() => {
        return !unref(groupFunc) ? (list!).filter(data => {
          const matcher = createMatcher && createMatcher(data);
          if (matcher && !matcher(filter)) return false;
          return true;
        }) : null
      })

      const childrenMatched = computed(() => {
        !unref(groupFunc) ? (list!).filter(data => {
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

      return {
        filtered,
        childrenMatched,
        listGrouped,
        ConnectionAppObject,
        SubDatabaseList
      }
    }
  })
</script>

<style scoped>

</style>
