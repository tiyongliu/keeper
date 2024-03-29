import {defineComponent, onMounted, PropType, ref, toRefs} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/second/keeper-tools'
import AppObjectList from '/@/second/appobj/AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {ConnectionsWithStatus, TablesNameSort} from '/@/second/typings/mysql'
import {useDatabaseList} from "/@/api/bridge";

export default defineComponent({
  name: "SubDatabaseList",
  props: {
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    data: {
      type: Object as PropType<ConnectionsWithStatus>
    },
    filter: {
      type: String as PropType<string>,
      default: ''
    }
  },
  setup(props) {
    const {data, filter, passProps} = toRefs(props)
    const databases = ref()

    onMounted(() => {
      useDatabaseList<TablesNameSort[]>({conid: data.value?._id}, databases)
    })

    return () => (
      <AppObjectList
        module={databaseAppObject}
        list={sortBy(
          (databases.value || []).filter(x => filterName(filter.value, x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({...db, connection: data.value})
        )}
        passProps={passProps.value}
      />
    )
  }
})
