import {defineComponent, PropType, toRefs, unref} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from 'keeper-tools'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {ConnectionsWithStatus, TablesNameSort} from '/@/second/typings/mysql'
// import {useDatabaseList} from "/@/api/metadataLoaders";
// import {metadataLoadersStore} from "/@/store/modules/metadataLoaders"
// import {useDatabaseList} from "/@/api/api";
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
    const databases = useDatabaseList<TablesNameSort[]>({conid: data.value?._id})

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
