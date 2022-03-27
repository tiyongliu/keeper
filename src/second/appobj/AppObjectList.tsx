import {computed, defineComponent, onMounted, PropType, unref} from 'vue'
import getLocalStorage from '/@/second/utility/getConnectionLabel'
import {filterName} from '/@/packages/tools/src/filterName'
import AppObjectListItem from './AppObjectListItem.vue'
export default defineComponent({
  name: "DatabaseWidget",
  components: {AppObjectListItem},
  props: {
    list: {
      type: Array as unknown as PropType<[]>,
    },
    groupFunc: {
      type: String as PropType<string>,
    }
  },
  setup(props) {
    // const dynamicList = computed(() => unref(props.list))
    const {list, groupFunc} = props

    console.log(list, '18')
    const filtered = computed(() => {
      !unref(groupFunc) ? list.filter(data => {

      }) : null
    })

    onMounted(() => {
      setTimeout(() => {
        console.log(props.list)
      }, 2000)
    })

    return () => (
      (list!).map(data => {
        return <AppObjectListItem {...data}/>
      })

    )
  }
})

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
