import {computed, defineComponent, onMounted, PropType, unref} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/packages/tools/src'
import FontIcon from '/@/second/icons/FontIcon.vue'
import './SubDatabaseList.less'
// import AppObjectList from './AppObjectList'
export default defineComponent({
  name: "SubDatabaseList",
  components: {
    // AppObjectList,
    FontIcon
  },
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
    },
  },
  setup(props) {
    const {data, filter} = props
    const databases = computed((): {name: string, sortOrder?: string}[] => {
      return [{"name":"information_schema"},{"name":"crmeb_java_beta"},{"name":"mysql"},{"name":"performance_schema"},{"name":"sys"}]
    })

    const isPinned = computed(() => {
      return true
    })

    console.log(computed(() => {
      return sortBy(
        (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
        x => x.sortOrder ?? x.name
      ).map(db => ({ ...db, connection: unref(data) }))
    }))

    onMounted(() => {

    })

    const onPin = (e) => {
      if (unref(isPinned)) {
        e?.stopPropagation()
        e?.preventDefault()
        console.log(data, `收藏`)
      }
    }

    return () => (

      <div class="main" draggable="true">
        <FontIcon icon="mdi mdi-database color-icon-gold" />crmeb_java_beta
        <span class="pin" >
           <FontIcon icon="mdi mdi-pin" onClick={onPin}/>
        </span>
      </div>

    )

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
