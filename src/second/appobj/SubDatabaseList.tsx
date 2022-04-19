import {computed, defineComponent, onMounted, PropType, unref, reactive} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/packages/tools/src'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
export default defineComponent({
  name: "SubDatabaseList",
  props: {
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
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
    const {data, filter, passProps} = props
    const databases = computed((): {name: string, sortOrder?: string}[] => {
      return [{"name":"information_schema"}, {"name":"crmeb_java_beta"},{"name":"mysql"},{"name":"performance_schema"},{"name":"sys"}]
    })

    console.log(computed(() => {
      return sortBy(
        (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
        x => x.sortOrder ?? x.name
      ).map(db => ({ ...db, connection: unref(data) }))
    }))

    onMounted(() => {
    })


    return () => (
      <AppObjectList
        module={databaseAppObject}
        list={sortBy(
          (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({ ...db, connection: unref(data) }))}
        passProps={unref(passProps)}
      />
    )
  }
})
