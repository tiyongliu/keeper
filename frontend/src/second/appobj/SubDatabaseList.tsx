import {defineComponent, onMounted, PropType, ref, toRefs, unref, watch} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from 'keeper-tools'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {ConnectionsWithStatus, TablesNameSort} from '/@/second/typings/mysql'
import {useDatabaseList} from "/@/api/metadataLoaders";

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
    const databases = ref<TablesNameSort[]>([])
    const showDatabases = async (conid: string) => {
      try {
        const result = await useDatabaseList({conid})
        if (Array.isArray(result)) {
          databases.value = result || []
        }
      } catch (e) {
        console.log(e)
        databases.value = []
      }
    }

    onMounted(() => {
      if (!!data.value?._id) {
        void showDatabases(data.value?._id)
      }
    })

    watch(() => data.value?._id, (conid) => {
      if (!!conid) {
        void showDatabases(conid!)
      }
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
