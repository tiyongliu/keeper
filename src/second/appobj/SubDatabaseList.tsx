import {computed, defineComponent, onMounted, PropType, unref} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/packages/tools/src'
// import AppObjectList from './AppObjectList'
export default defineComponent({
  name: "SubDatabaseList",
  // components: {AppObjectList},
  props: {
    passProps: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    data: {
      type: Object as PropType<{}>
    },
    filter: {
      type: String as PropType<string>,
      default: ''
    }
  },
  setup(props) {
    const {data, filter} = props
    const databases = computed((): {name: string, sortOrder?: string}[] => {
      return [{"name":"information_schema"},{"name":"crmeb_java_beta"},{"name":"mysql"},{"name":"performance_schema"},{"name":"sys"}]
    })




    console.log(computed(() => {
      return sortBy(
        (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
        x => x.sortOrder ?? x.name
      ).map(db => ({ ...db, connection: unref(data) }))
    }))

    onMounted(() => {

    })

    return () => (<div>1</div>)

    // return () => (
    //   <AppObjectList
    //     list={sortBy(
    //       (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
    //       x => x.sortOrder ?? x.name
    //     ).map(db => ({ ...db, connection: unref(data) }))}
    //   />
    // )
  }
})
