import {defineComponent, PropType, toRefs, unref, onMounted, ref} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/second/keeper-tools'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {ConnectionsWithStatus, TablesNameSort} from '/@/second/typings/mysql'
import {useDatabaseList} from "/@/api/sql";

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
          (unref(databases) || []).filter(x => filterName(filter.value, x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({...db, connection: data.value})
        )}
        passProps={unref(passProps)}
      />
    )
  }
})
