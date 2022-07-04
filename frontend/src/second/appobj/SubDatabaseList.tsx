import {defineComponent, onMounted, PropType, ref, toRefs, unref, watch} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from 'keeper-tools'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {ConnectionsWithStatus, TablesNameSort} from '/@/second/typings/mysql'
import {useDatabaseList} from "/@/api/metadataLoaders";
import {metadataLoadersStore} from "/@/store/modules/metadataLoaders"

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
    const metadataLoaders = metadataLoadersStore()

    const {data, filter, passProps} = toRefs(props)

    // const databases = metadataLoaders.getDatabaseList
    // const databases = ref<TablesNameSort[]>([])

    const showDatabases = async (conid: string) => {
      try {
        const result = await useDatabaseList(conid)
        console.log(`useDatabaseList`, result)
        if (Array.isArray(result)) {
          // databases.value = result
          metadataLoaders.subscribeDatabaseList(result)
        }
      } catch (e) {
        console.log(e)
      }
    }

    onMounted(() => {
      if (!!data.value?._id) void showDatabases(data.value?._id)
    })

    return () => (
      <AppObjectList
        module={databaseAppObject}
        list={sortBy(
          (unref(metadataLoaders.getDatabaseList) || []).filter(x => filterName(filter.value, x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({...db, connection: data.value})
        )}
        passProps={unref(passProps)}
      />
    )
  }
})
