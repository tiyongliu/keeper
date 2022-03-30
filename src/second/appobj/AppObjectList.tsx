import {computed, defineComponent, onMounted, PropType, unref} from 'vue'
import {compact} from 'lodash'
import AppObjectListItem from './AppObjectListItem.vue'

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
    // const dynamicList = computed(() => unref(props.list))
    const {list, groupFunc, filter, isExpandable} = props

    console.log(list, '18')
    const filtered = computed(() => {
      !unref(groupFunc) ? (list!).filter(data => {
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

    onMounted(() => {
    })

    return () => (
      (list!).map(data => {
        return <AppObjectListItem
          isHidden={!(filtered as unknown as []).includes(data)}
          data={data}
          isExpandable={isExpandable}
        />
      })

    )
  }
})

import getLocalStorage from '/@/second/utility/getConnectionLabel'
import {filterName} from '/@/packages/tools/src/filterName'

export const extractKey = data => data._id;
export const createMatcher = props => filter => {
  const { _id, displayName, server } = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(filter, displayName, server, ...databases.map(x => x.name));
};
export const createChildMatcher = props => filter => {
  if (!filter) return false;
  const { _id } = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(filter, ...databases.map(x => x.name));
};
