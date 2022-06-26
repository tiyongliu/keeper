import {defineComponent, unref, computed, PropType, toRefs} from 'vue'
import { isEqual, uniqWith, get } from 'lodash-es'
import AppObjectCore from './AppObjectCore.vue'
import { dataBaseStore } from "/@/store/modules/dataBase";
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d"
import {getDatabaseMenuItems} from './PinnedAppObject_'
export default defineComponent({
  name: 'DatabaseAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
    passProps: {
      type: Object as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
  },
  setup(props, {attrs}) {
    const dataBase = dataBaseStore()
    const {data, passProps} = toRefs(props)
    const isPinned = computed(() => {
      return dataBase.$state.pinnedDatabases.find(x => x.name == unref(data)!.name && x.connection?._id == unref(data)!.connection?._id)
    })

    const currentDatabase = computed(() => dataBase.getCurrentDatabase)

    return () => (
      <AppObjectCore
        {...attrs}
        title={unref(data)!.name}
        extInfo={unref(data)!.extInfo}
        icon="img database"
        isBold={get(unref(currentDatabase), 'connection._id') == get(unref(data)!.connection, '_id') &&
            get(unref(currentDatabase), 'name') == unref(data)!.name
        }
        onClick={() => dataBase.subscribeCurrentDatabase(unref(data)!)}
        menu={createMenu}
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        onPin={unref(isPinned) ? null : () => dataBase.subscribePinnedDatabases(uniqWith([
          ...unref(dataBase.$state.pinnedDatabases),
          unref(data!)
        ], isEqual))}
        onUnpin={unref(isPinned) ? () => {
          dataBase.subscribePinnedDatabases(
            dataBase.$state.pinnedDatabases.filter(x => x.name != unref(data)!.name || x.connection?._id != unref(data)!.connection?._id) as []
          )
        } : null}
      />
    )
  }
})

function createMenu() {
  return getDatabaseMenuItems()
}
