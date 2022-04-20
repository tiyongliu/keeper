import {computed, defineComponent, PropType, unref} from 'vue'
import {sortBy} from 'lodash-es'
import {filterName} from '/@/packages/tools/src'
import './SubDatabaseList.less'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
import {IPinnedDatabasesItem} from "/@/second/types/standard.d";

export default defineComponent({
  name: "SubDatabaseList",
  props: {
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
    data: {
      type: Object as PropType<IPinnedDatabasesItem>
    },
    filter: {
      type: String as PropType<string>,
      default: ''
    }
  },
  setup(props) {
    const {data, filter, passProps} = props
    const databases = computed((): { name: string, sortOrder?: string }[] => {
      return [{"name": "crmeb"}, {"name": "erd"}, {"name": "information_schema"}, {"name": "kb-dms"}, {"name": "mallplusbak"}, {"name": "mysql"}, {"name": "performance_schema"}, {"name": "schema"}, {"name": "shop_go"}, {"name": "sql_join"}, {"name": "ssodb"}, {"name": "yami_shops"}]
    })


    return () => (
      <AppObjectList
        module={databaseAppObject}
        list={sortBy(
          (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
            x => x.sortOrder ?? x.name
          ).map(db => ({...db, connection: unref(data)}) as unknown as IPinnedDatabasesItem
        )}
        passProps={unref(passProps)}
      />
    )
  }
})
